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
  no-cache-filter = ["generate"]
  output = ["./_site"]
}

target "vendor" {
  target = "vendor"
  output = ["."]
}

target "htmlproofer" {
  target = "htmlproofer"
  args = {
    JEKYLL_ENV = JEKYLL_ENV
  }
  output = ["type=cacheonly"]
}
