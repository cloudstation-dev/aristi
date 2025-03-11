/*
	Copyright 2025.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package controller

import (
	"context"
	"github.com/go-logr/logr"
	_ "istio.io/api/networking/v1alpha3"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"

	aristiv1alpha1 "cloudstation/aristi/api/v1alpha1"

	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// AristiReconciler reconciles a Aristi object
type AristiReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=aristi.cloudstation,resources=aristis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aristi.cloudstation,resources=aristis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aristi.cloudstation,resources=aristis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Aristi object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *AristiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("aristi", req.NamespacedName)
	// Obtener el CR Aristi

	var aristi aristiv1alpha1.Aristi
	if err := r.Get(ctx, req.NamespacedName, &aristi); err != nil {
		log.Error(err, "No se pudo obtener el recurso Aristi")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Construir VirtualService basado en el CRD
	found, err, result, err2 := createVirtualService(ctx, req, aristi, r, log)
	if err2 != nil {
		return result, err2
	}

	var canarySteps []argov1alpha1.CanaryStep

	for _, step := range aristi.Spec.Rollout.Strategy.Canary.Steps {
		if step.SetWeight != nil { // Solo agregar si tiene peso
			canarySteps = append(canarySteps, argov1alpha1.CanaryStep{
				SetWeight: step.SetWeight,
			})
		}

		if step.Pause != nil { // Solo agregar si tiene pausa
			canarySteps = append(canarySteps, argov1alpha1.CanaryStep{
				Pause: &argov1alpha1.RolloutPause{
					Duration: &intstr.IntOrString{
						Type:   1,
						IntVal: 0,
						StrVal: step.Pause.Duration.String(),
					},
				},
			})
		}
	}

	var containers []corev1.Container
	for _, c := range aristi.Spec.Rollout.Template.Spec.Containers {
		containers = append(containers, corev1.Container{
			Name:  c.Name,
			Image: c.Image,
		})
	}

	falcone := argov1alpha1.RolloutSpec{
		Replicas: aristi.Spec.Rollout.Replicas,
		Strategy: argov1alpha1.RolloutStrategy{
			Canary: &argov1alpha1.CanaryStrategy{
				CanaryService: aristi.Spec.Rollout.Strategy.Canary.CanaryService,
				StableService: aristi.Spec.Rollout.Strategy.Canary.StableService,
				TrafficRouting: &argov1alpha1.RolloutTrafficRouting{
					Istio: &argov1alpha1.IstioTrafficRouting{
						VirtualService: &argov1alpha1.IstioVirtualService{
							Name: aristi.Spec.Istio.VirtualService.Name,
							Routes: []string{
								"primary",
							},
						},
					},
				},
				Steps: canarySteps,
			},
		},
		Selector: aristi.Spec.Rollout.Selector,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: aristi.Spec.Rollout.Template.Labels,
			},
			Spec: corev1.PodSpec{
				Containers: containers,
			},
		},
	}

	rollout := &argov1alpha1.Rollout{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aristi.Name + "-rollout",
			Namespace: req.Namespace,
		},
		Spec: falcone,
	}

	// Aplicar Rollout en Kubernetes
	rolloutFound := &argov1alpha1.Rollout{}

	err = r.Get(ctx, client.ObjectKey{Name: rollout.Name, Namespace: rollout.Namespace}, found)
	if err != nil {
		log.Info("Creando Rollout", "name", rollout.Name)
		if err := r.Create(ctx, rollout); err != nil {
			log.Error(err, "No se pudo crear el Rollout")
			return ctrl.Result{}, err
		}
	} else {
		log.Info("Actualizando Rollout existente", "name", rollout.Name)
		rolloutFound.Spec = rollout.Spec

		if err := r.Update(ctx, found); err != nil {
			log.Error(err, "No se pudo actualizar el Rollout")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func createVirtualService(ctx context.Context, req ctrl.Request, aristi aristiv1alpha1.Aristi, r *AristiReconciler, log logr.Logger) (*istioclient.VirtualService, error, ctrl.Result, error) {
	var httpRoutes []*networkingv1alpha3.HTTPRoute

	httpRoute := &networkingv1alpha3.HTTPRoute{
		Name:  "primary",
		Route: []*networkingv1alpha3.HTTPRouteDestination{},
	}

	for _, route := range aristi.Spec.Istio.VirtualService.Routes {
		httpRoute.Route = append(httpRoute.Route, &networkingv1alpha3.HTTPRouteDestination{
			Destination: &networkingv1alpha3.Destination{
				Host: route.Destination.Host,
			},
			Weight: int32(route.Weight),
		})

		httpRoutes = append(httpRoutes, httpRoute)
	}

	virtualService := &istiov1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aristi.Spec.Istio.VirtualService.Name,
			Namespace: req.Namespace,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{"*"},
			Gateways: []string{"argo-gateway"},
			Http:     httpRoutes,
		},
	}

	// Apply VirtualService in Kubernetes
	virtualServicefound := &istiov1alpha3.VirtualService{}
	err := r.Get(ctx, client.ObjectKey{Name: virtualService.Name, Namespace: virtualService.Namespace}, virtualServicefound)
	if err != nil {
		log.Info("Creating VirtualService", "name", virtualService.Name)
		if err := r.Create(ctx, virtualService); err != nil {
			log.Error(err, "It couldn't create the VirtualService")
			return nil, nil, ctrl.Result{}, err
		}
	}

	log.Info("VirtualService created/updated correctly")
	return virtualServicefound, err, ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AristiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	_ = istioclient.AddToScheme(mgr.GetScheme())
	_ = argov1alpha1.AddToScheme(mgr.GetScheme())
	return ctrl.NewControllerManagedBy(mgr).
		For(&aristiv1alpha1.Aristi{}).
		Complete(r)
}
