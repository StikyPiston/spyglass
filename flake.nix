{
  description = "spyglass devshell and package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          name = "spyglass-devshell";

          packages = with pkgs; [
            go
            gopls
            gotools
            delve
          ];
        };

        packages.spyglass = pkgs.buildGoModule {
          pname = "spyglass";
          version = "2026.02.27-a";

          src = self;

          vendorHash = pkgs.lib.fakeHash;

          subPackages = [ "." ];
          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "An extensible search tool, inspired by Raycast and Vicinae";
            license = licenses.mit;
            platforms = platforms.linux;
          };
        };

        apps.spyglass = {
          type = "app";
          program = "${self.packages.${system}.spyglass}/bin/spyglass";
        };
      });
}
