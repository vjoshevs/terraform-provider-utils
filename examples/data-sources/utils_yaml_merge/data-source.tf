locals {
  yaml_1 = <<-EOT
    root:
      elem1: value1
      child1:
        cc1: 1
    list:
      - name: a1
        map:
          a1: 1
          b1: 1
      - name: a2
  EOT

  yaml_2 = <<-EOT
    root:
      elem2: value2
      child1:
        cc2: 2
    list:
      - name: a1
        map:
          a2: 2
      - name: a3
  EOT
}

data "utils_yaml_merge" "example" {
  input = [local.yaml_1, local.yaml_2]
}

output "output" {
  value = data.utils_yaml_merge.example.output
}

/* 
output = <<-EOT
  list:
      - map:
          a1: 1
          a2: 2
          b1: 1
        name: a1
      - name: a2
      - name: a3
  root:
      child1:
          cc1: 1
          cc2: 2
      elem1: value1
      elem2: value2
EOT
*/
