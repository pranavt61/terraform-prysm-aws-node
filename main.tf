
resource "null_resource" "this" {
  provisioner "local-exec" {
    command = "git clone https://github.com/stefa2k/ansible-prysm"
  }
}

