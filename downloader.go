package main

import (
	"archive/zip"
	"fmt"
	"github.com/bedrockmyshit/backend/app/minecraft"
	"github.com/joho/godotenv"
	"github.com/povsister/scp"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	args := os.Args
	dir := filepath.Dir(args[0])
	if len(args) < 3 {
		after, _ := strings.CutPrefix(args[0], dir+"\\")
		fmt.Printf("Usage: %s <ip:port> <ssh_addr>\n", after)
		return
	}
	serverAddr := args[1]

	split := strings.Split(serverAddr, ":")
	if len(split) != 2 {
		fmt.Println("Invalid ip:port")
		return
	}
	addr, port := split[0], split[1]

	username := os.Getenv("SSH_USERNAME")
	if username == "" {
		fmt.Println("SSH_USERNAME not set")
		return
	}
	password := os.Getenv("SSH_PASSWORD")
	if password == "" {
		fmt.Println("SSH_PASSWORD not set")
		return
	}

	fmt.Printf("[+] Connecting to %s\n", serverAddr)
	conn, err := minecraft.Connect(serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", serverAddr, err)
		return
	}
	defer conn.Close()

	packs := conn.ResourcePacks()
	if len(packs) == 0 {
		fmt.Printf("No packs found on %s\n", serverAddr)
		return
	}
	tempPath := os.Getenv("TEMP")

	f, err := os.Create(fmt.Sprintf("%s/%s.zip", tempPath, addr))
	if err != nil {
		_ = f.Close()
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	zipFile := zip.NewWriter(f)

	for _, pack := range packs {
		buf, err := minecraft.EncodePack(pack)
		if err != nil {
			fmt.Printf("Failed to encode pack: %v\n", err)
			continue
		}
		if pack.Encrypted() {
			buf, err = minecraft.DecryptPack(buf, pack.ContentKey())
			if err != nil {
				fmt.Printf("Failed to decrypt pack: %v\n", err)
				continue
			}
		}

		p, err := zipFile.Create(pack.Name() + ".zip")
		if err != nil {
			fmt.Printf("Failed to create zip file: %v\n", err)
			continue
		}
		if _, err = p.Write(buf); err != nil {
			fmt.Printf("Failed to write to zip file: %v\n", err)
			continue
		}
	}
	f.Close()
	zipFile.Close()
	fmt.Printf("[+] Packs saved to %s\n", f.Name())

	sshAddr := args[2]
	sshConf := scp.NewSSHConfigFromPassword(username, password)
	scpClient, err := scp.NewClient(sshAddr, sshConf, &scp.ClientOption{})
	defer scpClient.Close()
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", args[2], err)
		return
	}
	fmt.Printf("[+] Copying file %s to remote: %s\n", f.Name(), sshAddr)

	err = scpClient.CopyFileToRemote(f.Name(), fmt.Sprintf("/root/decryptmypack.com/packs/%s/%s/%s.zip", addr, port, addr), &scp.FileTransferOption{})
	if err != nil {
		fmt.Printf("Failed to copy file to remote: %v\n", err)
	}

	fmt.Println("[+] Done")
}
