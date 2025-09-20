{
  description = "Go development shell (pinned Go 1.24)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = [
          pkgs.go_1_24
          pkgs.gopls
        ];

        shellHook = ''
          echo "Go $(go version) ready to hack!"
        '';
      };
    };
}
