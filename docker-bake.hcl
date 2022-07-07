variable "JEKYLL_ENV" {
  default = "development"
}

target "_common" {
  args = {
    JEKYLL_ENV = JEKYLL_ENV
  }
  no-cache-filter = ["generate"]
}

group "default" {
  targets = ["release"]
}

target "release" {
  inherits = ["_common"]
  target = "release"
  output = ["./_site"]
}

target "vendor" {
  target = "vendor"
  output = ["."]
}

group "validate" {
  targets = ["htmlproofer", "mdl"]
}

target "htmlproofer" {
  inherits = ["_common"]
  target = "htmlproofer"
  output = ["type=cacheonly"]
}

target "mdl" {
  inherits = ["_common"]
  target = "mdl"
  output = ["type=cacheonly"]
}
