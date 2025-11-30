import 'dart:io';

class FileModel {
  final String name;
  final String path;
  final bool isDirectory;
  final int? size;
  final DateTime? modified;

  FileModel({
    required this.name,
    required this.path,
    required this.isDirectory,
    this.size,
    this.modified,
  });

  static FileModel fromFileSystemEntity(FileSystemEntity entity) {
    final stat = entity.statSync();
    final name = entity.path.split('/').last;

    return FileModel(
      name: name.isEmpty ? '/' : name,
      path: entity.path,
      isDirectory: stat.type == FileSystemEntityType.directory,
      size: stat.type == FileSystemEntityType.file ? stat.size : null,
      modified: stat.modified,
    );
  }

  String get displaySize {
    if (size == null) return '';

    if (size! < 1024) return '$size B';
    if (size! < 1024 * 1024) return '${(size! / 1024).toStringAsFixed(1)} KB';
    if (size! < 1024 * 1024 * 1024) return '${(size! / (1024 * 1024)).toStringAsFixed(1)} MB';
    return '${(size! / (1024 * 1024 * 1024)).toStringAsFixed(1)} GB';
  }

  String get displayModified {
    if (modified == null) return '';
    return '${modified!.day}/${modified!.month}/${modified!.year}';
  }
}