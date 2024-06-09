/* 
export ELEM1=value1
*/

locals {
  yaml_1 = <<-EOT
    root:
      elem1: !env ELEM1
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

output "output" {
  value = provider::utils::yaml_merge([local.yaml_1, local.yaml_2])
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