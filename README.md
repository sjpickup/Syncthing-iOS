# Syncthing-iOS

This is a Syncthing-iOS Port with Golang Libary for iOS. I have compiled it with Golang 1.4.2 and did some changes for iOS in Go and Syncthing.

1. Please install at first Golang https://golang.org/
2. Instal Xcode with the latest Version
3. Get a clone from Golang Solurce https://github.com/golang/go
4. Edit the file /misc/iso/clangwrap.sh 

    


   To build a cross compiling toolchain for iOS on OS X, first modify clangwrap.sh in misc/ios to match your setup. And then run:
    GOARM=7 CGO_ENABLED=1 GOARCH=arm CC_FOR_TARGET="pwd"/../misc/ios/clangwrap.sh \
     CXX_FOR_TARGET="pwd"/../misc/ios/clangwrap.sh ./make.bash
    To build a program, use the normal go build command:
    CGO_ENABLED=1 GOARCH=arm go build import/path
    
    Let’s just try those instructions he gave, without modifying anything:

   /* $ cd src
    $ GOARM=7 CGO_ENABLED=1 GOARCH=arm \
        CC_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh \
        CXX_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh \
        ./make.bash 
    ###### Building C bootstrap tool.
    cmd/dist
    ##### Building compilers and Go bootstrap tool for host, darwin/amd64.
    lib9
    libbio
    liblink
    cmd/gc
    ... */

    
    #Now you can compile Golang for iOS successful


5. Please change in /IDZWebBrowser/build-go.sh to your Project paths. (build-go.sh do all the job for you with gcc)
6. Inside of Xcode go to the tab "Build Phases"
    There create a Run Script when not existing and a our build-go.sh script
      1. Input Files
        $(SRCROOT)/main.go
      2. Output Files
        $(DERIVED_FILE_DIR)/glue-go.a
        

This is a very very basic how to, and troubles? write me... 


    




