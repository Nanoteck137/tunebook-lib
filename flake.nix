{
  description = "music library handling for tunebook";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    just.url = "github:casey/just/1.50.0";
    versionctl.url = "github:nanoteck137/versionctl/0.3.0";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        app = pkgs.buildGoModule {
          pname = "tunebook-lib";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/tunebook-lib"];

          ldflags = [
            "-X github.com/nanoteck137/tunebooklib.Version=${version}"
            "-X github.com/nanoteck137/tunebooklib.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "sha256-7K17JaXFsjf163g5PXCb5ng2gYdotnZ2IDKk8KFjNj0=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/tunebook-lib --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
          '';
        };
      in
      {
        packages = {
          default = app;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            imagemagick
            ffmpeg

            inputs.just.packages.${system}.default
            inputs.versionctl.packages.${system}.default
          ];
        };
      }
    );
}
