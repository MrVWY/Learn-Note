package 代码

import "fmt"

// 操作系统的文件系统是递归的： 文件夹中包含文件和文件夹， 其中又包含文件和文件夹， 以此类推。

// 每个文件和文件夹都可用一个 inode接口来表示。 inode接口中同样也有 clone克隆功能。

// file文件和 folder文件夹结构体都实现了 print打印和 clone方法， 因为它们都是 inode类型。 同时， 注意 file和 folder中的 clone方法。
// 这两者的 clone方法都会返回相应文件或文件夹的副本。 同时在克隆过程中， 我们会在其名称后面添加 “_clone” 字样。

func main() {
	file1 := &File{name: "File1"}
	file2 := &File{name: "File2"}
	file3 := &File{name: "File3"}

	folder1 := &Folder{
		children: []Inode{file1},
		name:     "Folder1",
	}

	folder2 := &Folder{
		children: []Inode{folder1, file2, file3},
		name:     "Folder2",
	}
	fmt.Println("\nPrinting hierarchy for Folder2")
	folder2.print("  ")

	cloneFolder := folder2.clone()
	fmt.Println("\nPrinting hierarchy for clone Folder")
	cloneFolder.print("  ")
}
