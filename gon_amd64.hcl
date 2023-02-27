source = ["./dist/darwin/macos-amd_darwin_amd64_v1/otterize"]
bundle_id = "com.otterize.cli.amd64"

apple_id {
  username = "ori@otterize.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Otterize, Inc."
}

zip {
  output_path = "./dist/darwin/macos-amd_darwin_amd64_v1/otterize_macOS_x86_64_notarized.zip"
}