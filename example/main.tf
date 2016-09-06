provider "foo" {}

resource "foo_res" "test" {
    must = "bar"
    option = "baz"
}