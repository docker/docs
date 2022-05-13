variable "JEKYLL_ENV" {
  default = "development"
}

group "default" {
  targets = ["release"]
}

target "release" {
  target = "release"
  args = {
    JEKYLL_ENV = JEKYLL_ENV
  }
  no-cache-filter = ["upstream-resources"]
  output = ["./_site"]
}

target "vendor" {
  target = "vendor"
  output = ["."]
}
