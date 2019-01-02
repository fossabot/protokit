# Infrastructure

...

## Prerequisties

* [Terraform](https://www.terraform.io/)
  * [Namecheap Terraform Provider](https://github.com/adamdecaf/terraform-provider-namecheap)
* [Google Cloud SDK](https://cloud.google.com/sdk/)
* [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Namecheap API Access](https://ap.www.namecheap.com/settings/tools/apiaccess/)

> __Heads up!__ Namecheap API Access is an opt-in feature. Once enabled, it's
> highly recommended you contact their support to improve response time.

### Configurations

The following configurations are expected as environment variables:

* `NAMECHEAP_TOKEN` - Namecheap API Access Token
* `NAMECHEAP_API_USER`- Namecheap API username
* `NAMECHEAP_USERNAME` - Namecheap account (probably the same as `NAMECHEAP_API_USER`)

## Getting started

> This bit is mostly note format for the time being\

### Initial setup

```sh
export GCP_PROJECT_ID=protokit
export GCP_REGION=us-west1
```

```sh
gcloud init
gcloud auth application-default login
gcloud config set project ${GCP_PROJECT_ID}
gcloud config unset container/use_client_certificate
```

#### Create the admin project

```sh
gcloud projects create ${GCP_PROJECT_ID} --set-as-default
```

> @NOTE: `terraform` definitions cannot interpolate variables, so I'm not sure
> if this bucket name can ever be configurable? Maybe an overrides? :thinking:

```sh
gsutil mb -p ${GCP_PROJECT_ID} -c regional -l ${GCP_REGION} gs://protokit-tf-state
gsutil versioning set on gs://protokit-tf-state
```

```
terraform init
terraform apply -auto-approve
```

## TODO

- [x] Manage domain nameservers
- [ ] Endpoint HTTPS
- [ ] E2E HTTPS (or up to the container at minimum)
