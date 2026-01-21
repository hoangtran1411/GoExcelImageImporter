# Roadmap and Future Improvements

This document records ideas and plans for upgrades to the **Go Excel Image Importer** application.

## 🚀 Priority Features (High Priority)

1.  **[x] Support More Formats**: Added support for `.webp` files.
2.  **Preview**: Allow previewing the list of product codes with missing images before running the process.
3.  **[x] Custom Image Size**: Allow users to input Excel cell dimensions or desired image size directly from the GUI.
4.  **[x] Logging**: Export a log file (`_missing.log`) for product codes that were not found.

## 🛠️ Technical Improvements

- [x] **Unit Tests**: Add more test cases for `internal/engine` (especially mapping and scaling logic).
- [x] **Concurrency Tuning**: Allow users to adjust the number of "Workers" from the interface to optimize for their machine's configuration.
- [x] **Wails Migration**: Migrated from Fyne to Wails for a more beautiful and lightweight interface.
- [x] **Auto Update**: Integrated automatic update mechanism via GitHub Releases.
- [x] **CI/CD**: Setup GitHub Actions to automatically build `.exe` files for new releases.

## 🎨 GUI Enhancements

- **[x] Dark Mode**: Modern, eye-friendly dark interface.
- **[x] Icons**: Custom application icon and professional look.
- **[x] Improved Progress Bar**: Displays detailed percentage and status.
- **[x] Responsive UI**: Interface automatically adapts to window size.

## 📚 Archived Documentation
- `ARCHITECTURE.md`: Detailed description of the Wails code structure.
- `.agent/skills/`: Reusable development skills, testing, and performance patterns.
