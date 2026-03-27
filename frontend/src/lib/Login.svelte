<script>
  import { login } from './api.js';

  let { appTitle = 'Kanban Board', onLogin } = $props();

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let submitting = $state(false);

  async function handleSubmit(e) {
    e.preventDefault();
    error = '';

    if (!email || !password) {
      error = 'Email and password are required.';
      return;
    }

    submitting = true;
    try {
      const user = await login(email, password);
      onLogin(user);
    } catch (err) {
      error = err.message;
    } finally {
      submitting = false;
    }
  }
</script>

<div class="login">
  <h1>{appTitle}</h1>

  <form onsubmit={handleSubmit}>
    <div class="field">
      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} required />
    </div>

    <div class="field">
      <label for="password">Password</label>
      <input id="password" type="password" bind:value={password} required />
    </div>

    {#if error}
      <p class="error">{error}</p>
    {/if}

    <button type="submit" disabled={submitting}>
      {submitting ? 'Signing in...' : 'Sign In'}
    </button>
  </form>
</div>

<style>
  .login {
    max-width: 360px;
    margin: 0 auto;
    padding: 48px 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: 100vh;
    justify-content: center;
  }

  h1 {
    font-size: 1.75rem;
    color: #333;
    margin: 0 0 32px;
  }

  form {
    width: 100%;
  }

  .field {
    margin-bottom: 16px;
  }

  label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: #555;
    margin-bottom: 4px;
  }

  input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 1rem;
    box-sizing: border-box;
  }

  input:focus {
    outline: none;
    border-color: #4a90d9;
    box-shadow: 0 0 0 2px rgba(74, 144, 217, 0.2);
  }

  .error {
    color: #c00;
    font-size: 0.875rem;
    margin: 12px 0;
  }

  button {
    width: 100%;
    padding: 10px;
    background: #4a90d9;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    margin-top: 8px;
  }

  button:hover:not(:disabled) {
    background: #357abd;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
