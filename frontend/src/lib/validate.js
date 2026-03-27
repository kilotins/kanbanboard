/**
 * Validate a password against the password policy:
 * - Minimum 8 characters
 * - At least one letter (uppercase or lowercase)
 * - At least one number
 *
 * @param {string} password
 * @returns {string} Error message, or empty string if valid
 */
export function validatePassword(password) {
  if (password.length < 8) {
    return 'Password must be at least 8 characters.';
  }
  if (!/[a-zA-Z]/.test(password)) {
    return 'Password must contain at least one letter.';
  }
  if (!/[0-9]/.test(password)) {
    return 'Password must contain at least one number.';
  }
  return '';
}
