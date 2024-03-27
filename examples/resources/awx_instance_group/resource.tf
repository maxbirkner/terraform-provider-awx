resource "awx_instance_group" "example" {
  name                       = "example"
  is_container_group         = true
  policy_instance_minimum    = 1
  policy_instance_percentage = 100
  pod_spec_override = jsonencode({
    spec = {
      containers = [
        {
          name = "awx_task",
          resources = {
            limits = {
              cpu    = "1000m"
              memory = "1Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "500Mi"
            }
          }
        }
      ]
    }
  })
}
