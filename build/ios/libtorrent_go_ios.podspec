Pod::Spec.new do |s|
    s.name             = 'libtorrent_go_ios'
    s.version          = '0.0.1'
    s.summary          = 'Torrent library for iOS'
    s.description      = <<-DESC
    Torrent library for iOS
                         DESC
    s.homepage         = 'http://github.com/yusing/go_torrent'
    s.license          = { :file => '../LICENSE' }
    s.author           = { 'Yusing' => 'yusing.wys@gmail.com' }
    s.source           = { :path => '.' }
    # s.public_header_files = 'Classes**/*.h'
    # s.source_files = '../../*.go'
    s.static_framework = true
    s.vendored_libraries = "*.a"
    s.dependency 'Flutter'
    s.platform = :ios, '11.0'
  
    # Flutter.framework does not contain a i386 slice.
    s.pod_target_xcconfig = { 'DEFINES_MODULE' => 'YES', 'EXCLUDED_ARCHS[sdk=iphonesimulator*]' => 'i386' }
    s.swift_version = '5.0'
  end