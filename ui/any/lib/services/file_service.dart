import 'dart:io';
import 'package:path_provider/path_provider.dart';
import '../models/file_model.dart';

class FileService {
  static Future<Directory> getRootDirectory() async {
    // Get the application documents directory as starting point
    return await getApplicationDocumentsDirectory();
  }

  static Future<List<FileModel>> getFilesInDirectory(String directoryPath) async {
    try {
      final directory = Directory(directoryPath);

      if (!await directory.exists()) {
        return [];
      }

      final entities = await directory.list().toList();

      // Sort: directories first, then files, both alphabetically
      entities.sort((a, b) {
        final aIsDir = a is Directory;
        final bIsDir = b is Directory;

        if (aIsDir != bIsDir) {
          return aIsDir ? -1 : 1;
        }

        return a.path.toLowerCase().compareTo(b.path.toLowerCase());
      });

      return entities
          .map((entity) => FileModel.fromFileSystemEntity(entity))
          .toList();
    } catch (e) {
      return [];
    }
  }

  static Future<bool> canNavigateToParent(String currentPath) async {
    try {
      final current = Directory(currentPath);
      final parent = current.parent;

      // Check if we can actually go to parent (not the same directory)
      return currentPath != parent.path;
    } catch (e) {
      return false;
    }
  }

  static String getParentPath(String currentPath) {
    final current = Directory(currentPath);
    return current.parent.path;
  }

  static Future<bool> createFolder(String parentPath, String folderName) async {
    try {
      final newFolderPath = '$parentPath/$folderName';
      final newFolder = Directory(newFolderPath);

      if (await newFolder.exists()) {
        return false; // Already exists
      }

      await newFolder.create(recursive: true);
      return true;
    } catch (e) {
      return false;
    }
  }

  static Future<List<FileModel>> searchFiles(String searchTerm) async {
    try {
      final rootDir = await getRootDirectory();
      final results = <FileModel>[];

      await _searchInDirectory(rootDir, searchTerm.toLowerCase(), results);

      return results;
    } catch (e) {
      return [];
    }
  }

  static Future<void> _searchInDirectory(
    Directory directory,
    String searchTerm,
    List<FileModel> results,
  ) async {
    try {
      final entities = await directory.list().toList();

      for (final entity in entities) {
        if (entity is File) {
          final fileName = entity.path.split('/').last.toLowerCase();
          if (fileName.contains(searchTerm)) {
            results.add(FileModel.fromFileSystemEntity(entity));
          }
        } else if (entity is Directory) {
          final dirName = entity.path.split('/').last.toLowerCase();
          if (dirName.contains(searchTerm)) {
            results.add(FileModel.fromFileSystemEntity(entity));
          }
          // Recursively search in subdirectories
          await _searchInDirectory(entity, searchTerm, results);
        }
      }
    } catch (e) {
      // Ignore directories we can't access
    }
  }
}