with import <nixpkgs> { }; {
  devEnv = stdenv.mkDerivation {
    name = "dev";
    buildInputs = [ stdenv go glibc.static ];
    CFLAGS = "-I${pkgs.glibc.dev}/include";
    LDFLAGS = "-L${pkgs.glibc}/lib";
  };
}
