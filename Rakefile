require 'rake/clean'
CLEAN.include FileList.new("tmp/*") # aruba's default tmp directory is local

desc "builds the binary"
task :build do
  # note starting in go1.18 this information will be available via go version -m
  version = `git describe --tags HEAD`
  sh "go", "build", "-o", "bin/scmpuff", "-mod=readonly", "-ldflags", "-X main.VERSION=#{version}"
end

desc "builds & installs the binary to $GOPATH/bin"
task :install => :build do
  # Don't cp directly over an existing file - it causes problems with Apple code signing.
  # https://developer.apple.com/documentation/security/updating_mac_software
  destination = "#{ENV['GOPATH']}/bin/scmpuff"
  rm destination if File.exist?(destination)
  cp "bin/scmpuff", destination
end

desc "run unit tests"
task :test do
  sh "go test ./..."
end

desc "run integration tests"
task :features => :build do
  sh "bundle exec cucumber -s --tags='not @wip'"
end

task "features:wip" => :build do
  sh "bundle exec cucumber -s --tags=@wip"
end

desc "package for distribution"
task :package do
  sh "goreleaser release --rm-dist --skip-publish"
end
CLOBBER.include "dist"

task :all => [:build, :test, :features]
task :default => :all
