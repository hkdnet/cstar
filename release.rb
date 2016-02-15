def root_dir
  File.expand_path(File.dirname(__FILE__))
end

def version
  cat = File.read("#{root_dir}/version.go")
  ver_line = cat.split("\n").find { |e| e.include?('Version') }
  ver_line.match(/"(.+)"/)[1]
end

def opt
  '-u hkdnet -r cstar'
end

def info(tag)
  `github-release info #{opt} -t #{tag}`
  $?.exitstatus
end

def create_release(tag)
  name = 'Stable Version'
  description = 'Colored stars generated from your projects'
  `github-release release #{opt} -t #{tag} -n "#{name}" --description "#{description}"`
end

def upload_file(tag)
  `github-release upload #{opt} -t #{tag} -n "cstar-osx.zip"     -f #{root_dir}/cstar-osx.zip`
  `github-release upload #{opt} -t #{tag} -n "cstar-win-x86.zip" -f #{root_dir}/cstar-win-x86.zip`
  `github-release upload #{opt} -t #{tag} -n "cstar-win-x64.zip" -f #{root_dir}/cstar-win-x64.zip`
end

def delete(tag)
  `github-release delete #{opt} -t #{tag}`
end

def upload(tag)
  delete(tag) if info(tag) == 0
  create_release(tag)
  upload_file(tag)
end

upload(version)
