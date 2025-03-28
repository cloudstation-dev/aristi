# Aristi – Stupidly Easy Progressive Delivery for Kubernetes 🚀🎯
Progressive delivery shouldn’t be complicated. That’s why we built Aristi, a Kubernetes operator designed to make progressive deployments effortless and insanely fast.

In its first version, Aristi simplifies canary deployments, enabling smooth, controlled rollouts in under a minute with just one manifest. But we’re not stopping there—Aristi is evolving to support Blue-Green deployments, A/B testing, and more, giving teams the flexibility to choose the best strategy for their needs.

Currently, Aristi integrates with ArgoCD and Istio, leveraging their power for seamless automation. Whether you're a startup or an enterprise, Aristi helps you deploy smarter, faster, and with zero hassle.

✨ *Progressive delivery should be simple. Aristi makes it happen.*

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Build and push your image

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/aristi:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

### To Deploy on the cluster

**Install the Custom Resource Definitions (CRDs) into the cluster:**

```sh
make install
```

**2. Option 1: Run the operator application locally**

If you'd like to test and debug the operator in your local environment, you can use the make run command. This compiles and runs the operator directly in your development environment, without the need to deploy it to the Kubernetes cluster.

``` sh
make run
```

**OR**
**Option 2: Deploy the operator to the cluster using the built image**

If you prefer to deploy the operator to the cluster, you can use the make deploy command. This will deploy the operator using the image specified by the IMG environment variable. Replace <some-registry>/aristi:tag with your image registry and tag.

```sh
make deploy IMG=<some-registry>/aristi:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**

Once the operator is deployed, you can create instances of your custom resources by applying the sample YAML files provided in the config/samples directory.

```sh
kubectl apply -f ./config/samples/aristi_v1alpha1_aristi.yaml
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/aristi:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/aristi/<tag or branch>/dist/install.yaml
```

## Contributing
Ever wanted to leave your mark on the cloud-native world? Well, now’s your chance! Aristi is on a mission to make Progressive Delivery for Kubernetes ridiculously fast, and we need your brilliance, your PRs, and maybe even your memes.

- 👨‍💻 Write code – because YAML doesn’t write itself (yet).
- 🔍 Report bugs – we promise not to blame you for finding them.
- 📖 Improve docs – help us make them clearer than your boss’s requirements.
- 🌍 Spread the word – because Aristi deserves more fame than a cat video.

Join us, contribute, and become a legend in the progressive deployment world! 🚀🐦

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

