import 'package:flutter/material.dart';
import '../models/file_model.dart';

class FileItem extends StatelessWidget {
  final FileModel file;
  final VoidCallback onTap;
  final VoidCallback? onLongPress;

  const FileItem({
    super.key,
    required this.file,
    required this.onTap,
    this.onLongPress,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(
        file.isDirectory ? Icons.folder : _getFileIcon(file.name),
        color: file.isDirectory ? Colors.blue : Colors.grey[600],
        size: 32,
      ),
      title: Text(
        file.name,
        style: const TextStyle(
          fontWeight: FontWeight.w500,
          fontSize: 16,
        ),
      ),
      subtitle: !file.isDirectory && file.size != null
          ? Text(
              '${file.displaySize} • ${file.displayModified}',
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[600],
              ),
            )
          : file.modified != null
              ? Text(
                  file.displayModified,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                  ),
                )
              : null,
      trailing: file.isDirectory
          ? Icon(Icons.chevron_right, color: Colors.grey[400])
          : null,
      onTap: onTap,
      onLongPress: onLongPress,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
    );
  }

  IconData _getFileIcon(String fileName) {
    final extension = fileName.split('.').last.toLowerCase();

    switch (extension) {
      case 'pdf':
        return Icons.picture_as_pdf;
      case 'doc':
      case 'docx':
        return Icons.description;
      case 'xls':
      case 'xlsx':
        return Icons.table_chart;
      case 'ppt':
      case 'pptx':
        return Icons.slideshow;
      case 'txt':
      case 'md':
        return Icons.text_snippet;
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
      case 'bmp':
        return Icons.image;
      case 'mp3':
      case 'wav':
      case 'flac':
        return Icons.audiotrack;
      case 'mp4':
      case 'avi':
      case 'mov':
        return Icons.video_file;
      case 'zip':
      case 'rar':
      case '7z':
      case 'tar':
      case 'gz':
        return Icons.archive;
      case 'html':
      case 'css':
      case 'js':
      case 'dart':
      case 'py':
      case 'java':
      case 'cpp':
      case 'c':
        return Icons.code;
      default:
        return Icons.insert_drive_file;
    }
  }
}