import 'package:flutter/material.dart';
import '../models/file_model.dart';
import '../services/file_service.dart';
import 'file_item.dart';

class FileBrowser extends StatefulWidget {
  const FileBrowser({super.key});

  @override
  State<FileBrowser> createState() => _FileBrowserState();
}

class _FileBrowserState extends State<FileBrowser> {
  List<FileModel> _files = [];
  final List<String> _pathHistory = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadFiles();
  }

  Future<void> _loadFiles() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final currentPath = _pathHistory.isEmpty
          ? (await FileService.getRootDirectory()).path
          : _pathHistory.last;

      final files = await FileService.getFilesInDirectory(currentPath);

      // Add parent directory navigation if not at root
      if (await FileService.canNavigateToParent(currentPath)) {
        final parentPath = FileService.getParentPath(currentPath);
        files.insert(
          0,
          FileModel(
            name: '..',
            path: parentPath,
            isDirectory: true,
          ),
        );
      }

      setState(() {
        _files = files;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Failed to load files: $e';
        _isLoading = false;
      });
    }
  }

  void _navigateToDirectory(FileModel file) {
    if (file.name == '..') {
      _navigateToParent();
    } else if (file.isDirectory) {
      setState(() {
        _pathHistory.add(file.path);
      });
      _loadFiles();
    } else {
      _showFileDetails(file);
    }
  }

  void _navigateToParent() {
    if (_pathHistory.length > 1) {
      setState(() {
        _pathHistory.removeLast();
      });
      _loadFiles();
    } else {
      // Navigate to root
      setState(() {
        _pathHistory.clear();
      });
      _loadFiles();
    }
  }

  void _showFileDetails(FileModel file) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(file.name),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Path: ${file.path}'),
            if (file.size != null) Text('Size: ${file.displaySize}'),
            if (file.modified != null) Text('Modified: ${file.displayModified}'),
            Text('Type: ${file.isDirectory ? 'Directory' : 'File'}'),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  Future<void> _showCreateFolderDialog() async {
    final controller = TextEditingController();
    final currentPath = _pathHistory.isEmpty
        ? (await FileService.getRootDirectory()).path
        : _pathHistory.last;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create New Folder'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            labelText: 'Folder Name',
            hintText: 'Enter folder name',
          ),
          autofocus: true,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              final folderName = controller.text.trim();
              if (folderName.isNotEmpty) {
                final success = await FileService.createFolder(
                  currentPath,
                  folderName,
                );
                if (success) {
                  Navigator.pop(context);
                  _loadFiles();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Folder "$folderName" created')),
                  );
                } else {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Failed to create folder')),
                  );
                }
              }
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }

  Future<void> _showSearchDialog() async {
    final controller = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Search Files'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            labelText: 'Search',
            hintText: 'Enter file or folder name',
          ),
          autofocus: true,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              final searchTerm = controller.text.trim();
              if (searchTerm.isNotEmpty) {
                Navigator.pop(context);
                _performSearch(searchTerm);
              }
            },
            child: const Text('Search'),
          ),
        ],
      ),
    );
  }

  Future<void> _performSearch(String searchTerm) async {
    setState(() {
      _isLoading = true;
    });

    final searchResults = await FileService.searchFiles(searchTerm);

    setState(() {
      _files = searchResults;
      _isLoading = false;
    });

    if (searchResults.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('No results found for "$searchTerm"')),
      );
    }
  }

  String _getCurrentPath() {
    if (_pathHistory.isEmpty) {
      return '/ Root';
    }
    return _pathHistory.last;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'File Browser',
          style: TextStyle(fontSize: 20),
        ),
        backgroundColor: Colors.blue[600],
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: _showSearchDialog,
            tooltip: 'Search',
          ),
          IconButton(
            icon: const Icon(Icons.home),
            onPressed: () {
              setState(() {
                _pathHistory.clear();
              });
              _loadFiles();
            },
            tooltip: 'Home',
          ),
        ],
      ),
      body: Column(
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            color: Colors.grey[100],
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Current Path:',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  _getCurrentPath(),
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Colors.grey[800],
                  ),
                ),
              ],
            ),
          ),
          if (_isLoading)
            const Expanded(
              child: Center(
                child: CircularProgressIndicator(),
              ),
            )
          else if (_error != null)
            Expanded(
              child: Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.error, size: 64, color: Colors.red[400]),
                    const SizedBox(height: 16),
                    const Text(
                      'Error loading files',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(_error!),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: _loadFiles,
                      child: const Text('Retry'),
                    ),
                  ],
                ),
              ),
            )
          else if (_files.isEmpty)
            Expanded(
              child: Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.folder_open,
                        size: 64, color: Colors.grey[400]),
                    const SizedBox(height: 16),
                    Text(
                      'No files found',
                      style: TextStyle(
                        fontSize: 18,
                        color: Colors.grey[600],
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'This directory is empty',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.grey[500],
                      ),
                    ),
                  ],
                ),
              ),
            )
          else
            Expanded(
              child: ListView.builder(
                itemCount: _files.length,
                itemBuilder: (context, index) {
                  final file = _files[index];
                  return FileItem(
                    file: file,
                    onTap: () => _navigateToDirectory(file),
                  );
                },
              ),
            ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateFolderDialog,
        backgroundColor: Colors.blue[600],
        tooltip: 'Add Folder',
        child: const Icon(Icons.add),
      ),
    );
  }
}