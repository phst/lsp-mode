;;; lsp-mode-tests.el --- unit tests for lsp-mode.el  -*- lexical-binding: t; -*-

;; Copyright (C) 2017  Google Inc.

;; Author: Philipp Stephani <phst@google.com>

;; This program is free software; you can redistribute it and/or modify
;; it under the terms of the GNU General Public License as published by
;; the Free Software Foundation, either version 3 of the License, or
;; (at your option) any later version.

;; This program is distributed in the hope that it will be useful,
;; but WITHOUT ANY WARRANTY; without even the implied warranty of
;; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
;; GNU General Public License for more details.

;; You should have received a copy of the GNU General Public License
;; along with this program.  If not, see <http://www.gnu.org/licenses/>.

;;; Commentary:

;; Unit tests for lsp-mode.el.

;;; Code:

(require 'lsp-mode)

(require 'ert)
(require 'f)

(defconst lsp-mode-tests--directory
  (file-name-directory (f-this-file))
  "Directory where this file resides in.")

(lsp-define-stdio-client
 test-client "go" (lambda () lsp-mode-tests--directory) nil
 :command-fn #'lsp-mode-test--command)

(ert-deftest lsp-define-stdio-client ()
  (should (fboundp 'test-client-enable)))

(defvar lsp-mode-test--server-flags)

(ert-deftest lsp-mode/server-crashes ()
  (with-temp-buffer
    (let ((state-file (make-temp-file "test-server-state.json")))
      (unwind-protect
          (let ((debug-on-error t)
                (buffer-file-name
                 (expand-file-name "test.go" lsp-mode-tests--directory))
                (lsp-mode-test--server-flags
                 (list "-state_file" state-file "-initial_crashes" "1")))
            (test-client-enable))
        (delete-file state-file)
        (with-current-buffer "*test-client stderr*"
          (princ "Standard error from test server:")
          (terpri)
          (princ (buffer-string))
          (terpri))))))

(defun lsp-mode-test--command ()
  (cons (expand-file-name "test_server" lsp-mode-tests--directory)
        lsp-mode-test--server-flags))

;;; lsp-mode-tests.el ends here
