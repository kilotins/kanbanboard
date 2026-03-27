<script>
  import { postSetup } from './api.js';
  import { validatePassword } from './validate.js';

  let { onComplete } = $props();

  let appTitle = $state('Kanban Board');
  let name = $state('');
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let error = $state('');
  let submitting = $state(false);

  async function handleSubmit(e) {
    e.preventDefault();
    error = '';

    if (!name || !email || !password) {
      error = 'All fields are required.';
      return;
    }

    if (password !== confirmPassword) {
      error = 'Passwords do not match.';
      return;
    }

    const passwordError = validatePassword(password);
    if (passwordError) {
      error = passwordError;
      return;
    }

    submitting = true;
    try {
      await postSetup({ name, email, password, appTitle });
      onComplete();
    } catch (err) {
      error = err.message;
    } finally {
      submitting = false;
    }
  }
</script>

<div class="onboarding">
  <h1>Welcome</h1>
  <p class="subtitle">Set up your Kanban Board</p>

  <form onsubmit={handleSubmit}>
    <div class="field">
      <label for="appTitle">Application Title</label>
      <input id="appTitle" type="text" bind:value={appTitle} placeholder="Kanban Board" />
    </div>

    <hr />

    <h2>Administrator Account</h2>

    <div class="field">
      <label for="name">Name</label>
      <input id="name" type="text" bind:value={name} required />
    </div>

    <div class="field">
      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} required />
    </div>

    <div class="field">
      <label for="password">Password</label>
      <input id="password" type="password" bind:value={password} required />
    </div>

    <div class="field">
      <label for="confirmPassword">Confirm Password</label>
      <input id="confirmPassword" type="password" bind:value={confirmPassword} required />
    </div>

    {#if error}
      <p class="error">{error}</p>
    {/if}

    <button type="submit" disabled={submitting}>
      {submitting ? 'Setting up...' : 'Complete Setup'}
    </button>
  </form>
</div>

<style>
  .onboarding {
    max-width: 420px;
    margin: 0 auto;
    padding: 48px 24px;
  }

  h1 {
    font-size: 1.75rem;
    margin: 0 0 4px;
    color: #333;
  }

  h2 {
    font-size: 1.1rem;
    margin: 0 0 12px;
    color: #555;
  }

  .subtitle {
    color: #666;
    margin: 0 0 32px;
  }

  hr {
    border: none;
    border-top: 1px solid #e0e0e0;
    margin: 24px 0;
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
