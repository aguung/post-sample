data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/loader"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://pgvector/pgvector/pg17/dev"
  migration {
    dir = "file://migrations"
    revisions_schema = "public"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
