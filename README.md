# Syncthing-iOS

This is a Syncthing-iOS Port with Golang Libary for iOS. I have compiled it with Golang 1.4.2 and did some changes for iOS in Go and Syncthing.

1. Please install at first Golang https://golang.org/
2. Instal Xcode with the latest Version
3. Get a clone from Golang Solurce https://github.com/golang/go
4. Edit the file /misc/iso/clangwrap.sh 
    Change this lines to your settings
    
    SDK=your sdk version
    SDK_PATH=`xcrun --sdk $SDK --show-sdk-path`
5. Compile Golang from Source (You need to do the 1 Step, golang is requered to be installed before you can compile from source)
6. Please change in /IDZWebBrowser/build-go.sh to your Project paths.
7. Inside of Xcode go to the tab "Build Phases"
    There create a Run Script when not existing and a our build-go.sh script
      1. Input Files
        $(SRCROOT)/main.go
      2. Output Files
        $(DERIVED_FILE_DIR)/glue-go.a
        

This is a very very basic how to, and troubles? write me... 


    




