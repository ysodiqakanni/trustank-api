## To create a new helm chart
run `helm create project-name`

### To perform a dry run after changing the chart,
Inside the helm folder, run `helm install --dry-run my-release trustank-api`

**To deploy the chart**
run `helm install backend trustank-api`  # where 'backend' is the release name