require 'rake/clean'

# aruba's default tmp directory is local
CLEAN.include FileList.new("tmp/*")

# runs the generate script, which will bootstrap anything it needs in script
desc "generates bindata files"
task :generate do
  sh "go generate ./..."
end

# the unix build script does not force `generate` as prereq, but the task here
# does since we want to always make sure to be up to date with any changes made
desc "builds the binary"
task :build => :generate do
  sh "script/build"
end

desc "builds & installs the binary to $GOPATH/bin"
task :install => :build do
  cp "bin/scmpuff", "#{ENV['GOPATH']}/bin/scmpuff"
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
