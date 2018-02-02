;;; test-helper.el --- Helpers for lsp-mode-test.el  -*- lexical-binding: t; -*-

;;; Commentary:

;; Initializes test support for ‘lsp-mode’.

;;; Code:

(require 'f)

(add-to-list 'load-path
             (file-name-as-directory (f-parent (f-parent (f-this-file)))))

;;; test-helper.el ends here
