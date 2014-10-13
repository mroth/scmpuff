desc "build the binary"
task :build do
  sh "script/build"
end

desc "run unit tests"
task :test do
  sh "go test ./..."
end

desc "run feature tests"
task :features => :build do
  sh "cucumber -s --tags=~@wip"
end

task :all => [:build, :test, :features]
