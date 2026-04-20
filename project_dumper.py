#!/usr/bin/env python3
"""
Скрипт для сбора содержимого всех файлов проекта в один файл.
Использование: python dump_project.py [output_file]
По умолчанию результат записывается в project_dump.txt
"""

import os
import sys
from pathlib import Path

# Директории, которые следует игнорировать
IGNORED_DIRS = {
    '.git', '.idea', '.vscode', '__pycache__', 'node_modules',
    'vendor', 'dist', 'build', 'target', '.gradle', '.mvn',
    'venv', '.venv', 'env', '.env', 'bin', 'obj',
    '.next', '.nuxt', 'coverage', '.pytest_cache', '.mypy_cache',
    '.tox', 'htmlcov', '.cache', 'tmp', 'temp',
}

# Расширения файлов, которые следует игнорировать (бинарные и сгенерированные)
IGNORED_EXTENSIONS = {
    # Бинарные
    '.exe', '.dll', '.so', '.dylib', '.a', '.o', '.obj', '.lib',
    '.class', '.jar', '.war', '.ear', '.pyc', '.pyo',
    # Изображения
    '.png', '.jpg', '.jpeg', '.gif', '.bmp', '.ico', '.svg', '.webp',
    '.tiff', '.tif', '.raw', '.psd',
    # Видео и аудио
    '.mp4', '.avi', '.mov', '.mkv', '.webm', '.mp3', '.wav', '.flac', '.ogg',
    # Архивы
    '.zip', '.tar', '.gz', '.bz2', '.xz', '.7z', '.rar',
    # Документы
    '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx',
    # Шрифты
    '.ttf', '.otf', '.woff', '.woff2', '.eot',
    # Прочее
    '.db', '.sqlite', '.sqlite3', '.lock',
    '.DS_Store',
}

# Конкретные файлы, которые следует игнорировать
IGNORED_FILES = {
    '.DS_Store', 'Thumbs.db', 'go.sum', 'package-lock.json',
    'yarn.lock', 'pnpm-lock.yaml', 'Cargo.lock', 'poetry.lock',
}

# Максимальный размер файла в байтах (1 МБ)
MAX_FILE_SIZE = 1024 * 1024


def should_skip_dir(dirname: str) -> bool:
    return dirname in IGNORED_DIRS or dirname.startswith('.')


def should_skip_file(filepath: Path) -> tuple[bool, str]:
    if filepath.name in IGNORED_FILES:
        return True, "ignored filename"
    if filepath.suffix.lower() in IGNORED_EXTENSIONS:
        return True, f"ignored extension ({filepath.suffix})"
    try:
        size = filepath.stat().st_size
        if size > MAX_FILE_SIZE:
            return True, f"too large ({size} bytes)"
    except OSError as e:
        return True, f"stat error: {e}"
    return False, ""


def is_binary(filepath: Path) -> bool:
    """Проверка, является ли файл бинарным (по наличию нулевых байтов)."""
    try:
        with open(filepath, 'rb') as f:
            chunk = f.read(8192)
            return b'\x00' in chunk
    except OSError:
        return True


def dump_project(root: Path, output: Path) -> None:
    total_files = 0
    skipped_files = 0
    total_bytes = 0

    with open(output, 'w', encoding='utf-8') as out:
        # Заголовок с деревом проекта
        out.write("=" * 80 + "\n")
        out.write(f"PROJECT DUMP: {root.resolve()}\n")
        out.write("=" * 80 + "\n\n")

        # Сбор всех файлов
        all_files = []
        for dirpath, dirnames, filenames in os.walk(root):
            # Фильтрация директорий на месте
            dirnames[:] = [d for d in dirnames if not should_skip_dir(d)]
            for filename in filenames:
                all_files.append(Path(dirpath) / filename)

        all_files.sort()

        # Сначала вывод дерева файлов
        out.write("FILE TREE:\n")
        out.write("-" * 80 + "\n")
        for f in all_files:
            rel = f.relative_to(root)
            skip, reason = should_skip_file(f)
            marker = f"  [SKIP: {reason}]" if skip else ""
            out.write(f"  {rel}{marker}\n")
        out.write("\n" + "=" * 80 + "\n\n")

        # Затем содержимое файлов
        for filepath in all_files:
            rel = filepath.relative_to(root)
            skip, reason = should_skip_file(filepath)
            if skip:
                skipped_files += 1
                continue

            if is_binary(filepath):
                skipped_files += 1
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE: {rel} [BINARY - SKIPPED]\n")
                out.write(f"{'=' * 80}\n")
                continue

            try:
                with open(filepath, 'r', encoding='utf-8', errors='replace') as f:
                    content = f.read()
            except OSError as e:
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE: {rel} [ERROR: {e}]\n")
                out.write(f"{'=' * 80}\n")
                skipped_files += 1
                continue

            line_count = content.count('\n') + 1
            out.write(f"\n{'=' * 80}\n")
            out.write(f"FILE: {rel}\n")
            out.write(f"LINES: {line_count} | SIZE: {len(content)} bytes\n")
            out.write(f"{'=' * 80}\n")
            out.write(content)
            if not content.endswith('\n'):
                out.write('\n')

            total_files += 1
            total_bytes += len(content)

        # Итоговая статистика
        out.write(f"\n{'=' * 80}\n")
        out.write("SUMMARY\n")
        out.write(f"{'=' * 80}\n")
        out.write(f"Total files dumped: {total_files}\n")
        out.write(f"Total files skipped: {skipped_files}\n")
        out.write(f"Total bytes: {total_bytes}\n")

    print(f"✓ Dump saved to: {output}")
    print(f"  Files dumped:  {total_files}")
    print(f"  Files skipped: {skipped_files}")
    print(f"  Total size:    {total_bytes} bytes")


def main() -> None:
    root = Path('.')
    output = Path(sys.argv[1]) if len(sys.argv) > 1 else Path('project_dump.txt')
    dump_project(root, output)


if __name__ == '__main__':
    main()