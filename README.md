# uiharuam

Uiharu Archive Manager. Basically focusing on tapes.

## Usage

It's based on **tar**, but instead of writing files to tapes directly, it first generates **filelists** for tar, and has a 'write' subcommand that acts as a helper to perform the actual write. If you prefer, you can also write manually using the tar command with generated filelists. (But don't forget to create a "WRITTEN_FILELISTS" file to save your progress, and there must be a '\n' before the EOF)

Not like traditional tar-style multi-volume archives, it calculates the file size first, and generates file lists for creating separate, single-volume tars for each tape. This means that if you have a file 'X' on tape n, you don't need to start reading from tape 1.

Another advantage of it is it will automatically generates a meta database, including:

- Path of file or directory
- SHA-512 checksum
- File size
- Position (which tape, which tar)

You can use any SQLite tool to read or query.

Appending files and finding files are not supported at the moment, but these may come in the future.

This is a libre software under **AGPLv3**, and comes with **ABSOLUTELY NO WARRANTY**, to the extent permitted by applicable law.
