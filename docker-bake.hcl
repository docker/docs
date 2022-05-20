variable "JEKYLL_ENV" {
  default = "development"
}
variable "DOCS_URL" {
  default = "http://localhost:4000"
}

target "_common" {
  args = {
    JEKYLL_ENV = JEKYLL_ENV
    DOCS_URL = DOCS_URL
  }
}

group "default" {
  targets = ["release"]
}

target "release" {
  inherits = ["_common"]
  target = "release"
  no-cache-filter = ["generate"]
  output = ["./_site"]
}

target "vendor" {
  target = "vendor"
  output = ["."]
}

target "htmlproofer" {
  inherits = ["_common"]
  target = "htmlproofer"
  output = ["type=cacheonly"]
}
