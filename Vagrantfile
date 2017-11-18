# -*- mode: ruby -*-
# # vi: set ft=ruby :

# origins:

# http://blog.scottlowe.org/2014/10/22/multi-machine-vagrant-with-yaml/
require 'yaml'
require 'pp'

basedir = ENV.fetch('USERPROFILE', '')
basedir = ENV.fetch('HOME', '') if basedir == ''
basedir = basedir.gsub('\\', '/')
debug = ENV.fetch('DEBUG', '')
debug = ( debug =~ (/^(true|yes|1)$/i))

dir = File.expand_path(File.dirname(__FILE__))

# Read nodes details
nodes = {}
nodes = YAML.load(<<-NODES)
---
- name: 'uru'
  box: 'centos'
  ram: '512'
  ipaddress: '172.17.8.102'

NODES
if debug
  nodes.each do |box|
    pp box
  end
end

# Read box details
configs = YAML::load( <<-BOXES)
---
:boot: 'centos'
# centos 7.3 x64 with Puppet 4.10.x
# origin:
# https://app.vagrantup.com/mbrush/boxes/centos7-puppet/versions/1.1.1/providers/virtualbox.box
'centos':
  :image_name: 'centos'
  :box_memory: '512'
  :box_cpus: '1'
  :box_gui: false
  :config_vm_newbox: false
  :config_vm_default: 'linux'
  :config_vm_box: 'centos'
  :image_filename: 'centos7-puppet-x86_64.box'

BOXES
if debug
  pp configs
end
modulepath = '/opt/puppetlabs/puppet/modules'

$puppet_prereq_script = <<-SCRIPT
#!/usr/bin/env bash

# force the locale change
cat<<EOF>/etc/environment
LANG=en_US.utf-8
LC_ALL=en_US.utf-8
EOF
MODULE_PATH='#{modulepath}'

if [ ! -d ${MODULE_PATH}/stdlib ]; then
  puppet module install 'puppetlabs-stdlib' --version '4.22.0'
fi
if [ ! -d ${MODULE_PATH}/java ]; then
  puppet module install 'puppetlabs-java' --version '1.3.0'
fi
# for encrypted 
# gem install hiera-eyaml --no-rdoc --no-ri
SCRIPT

# bind_ip=127.0.0.1

box_config = {}

Vagrant.configure('2') do |config|
  nodes.each do |box|
    box_name = box['name']
    box_config = configs[box['box']]
    box_gui = box_config[:box_gui] != nil && box_config[:box_gui].to_s.match(/(true|yes|1)$/i) != nil
    box_cpus = box_config[:box_cpus].to_i
    box_memory = box_config[:box_memory].to_i
    newbox = box_config[:config_vm_newbox]
    image_filename = box_config[:image_filename]
    box_url = "file://#{basedir}/Downloads/#{image_filename}"
    config.vm.define box_name do |guest|
      guest.vm.box = box_config[:image_name]
      guest.vm.box_url = box_url
      guest.vm.network 'private_network', ip: box['ipaddress']
      guest.vm.provider :virtualbox do |vb|
        vb.name = box_name
        vb.memory = box_memory
      end

      config.vm.synced_folder './' , '/vagrant'
      config.vm.provision 'shell', inline: $puppet_prereq_script

      manifest = 'linux_uru.pp'

      # workaround for older Vagrant / Puppet 4 compatibility issue
      # NOTE: hiera arguments not set
      # config.vm.provision 'shell', inline: "puppet apply --modulepath=#{modulepath}:/vagrant/modules /vagrant/manifests/#{manifest}"

       config.vm.provision :puppet do |puppet|
          puppet.hiera_config_path = 'hiera.yaml'
          puppet.module_path    = 'modules'
          puppet.manifests_path = 'manifests'
          puppet.manifest_file  = manifest
          puppet.options        = "--verbose --modulepath #{modulepath}:/vagrant/modules "
	  # hack to have hiera under Vagrant. des not seem to work
          puppet.working_directory = '/tmp/vagrant-puppet'
          config.vm.synced_folder 'hiera/', '/tmp/vagrant-puppet/hiera'
          config.vm.synced_folder 'keys/', '/tmp/vagrant-puppet/keys'
       end
    end
  end
end
