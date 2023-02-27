source = ["./dist/darwin/macos-arm_darwin_arm64/otterize"]
bundle_id = "com.otterize.cli.arm64"

apple_id {
  username = "ori@otterize.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Otterize, Inc."
}

zip {
  output_path = "./dist/darwin/macos-arm_darwin_arm64/otterize_macOS_arm64_notarized.zip"
}