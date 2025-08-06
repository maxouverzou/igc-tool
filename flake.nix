{
  description = "Development environment for igc-tool";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";

  outputs = {self, ...}@inputs:
    let
      goVersion = 24; # Change this to update the whole stack

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: inputs.nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import inputs.nixpkgs {
          inherit system;
          overlays = [ inputs.self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      packages = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.buildGoModule {
          pname = "igc-tool";
          version = "0.1.1";
          src = self.outPath;
          vendorHash = "sha256-SR3LZYQ1bWCwgYt9TIBt0YazYUWU+2jUp8x056SOcdU=";
          
          ldflags = [
            "-s" "-w"
            "-X igc-tool/internal/version.Version=0.1.1"
            "-X igc-tool/internal/version.GitCommit=${self.rev or "dirty"}"
            "-X igc-tool/internal/version.BuildDate=${pkgs.lib.optionalString (self ? lastModified) (builtins.substring 0 8 (builtins.toString self.lastModified))}"
          ];
          
          meta = with pkgs.lib; {
            description = "IGC track analysis tools";
            homepage = "https://github.com/maxouverzou/igc-track-tools";
            license = licenses.mit; # Adjust license as needed
            maintainers = [ ];
          };
        };
      });

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            # go (version is specified by overlay)
            go

            # goimports, godoc, etc.
            gotools

            # https://github.com/golangci/golangci-lint
            golangci-lint
            
            inputs.self.packages.${pkgs.system}.default
          ];
        };
      });
    };
}
