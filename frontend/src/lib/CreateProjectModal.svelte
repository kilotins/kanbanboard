<script>
  import { createProject } from './api.js';

  let { onCreated, onCancel } = $props();

  let name = $state('');
  let error = $state('');
  let submitting = $state(false);

  async function handleSubmit(e) {
    e.preventDefault();
    error = '';

    if (!name.trim()) {
      error = 'Project name is required.';
      return;
    }

    submitting = true;
    try {
      const project = await createProject(name.trim());
      onCreated(project);
    } catch (err) {
      error = err.message;
    } finally {
      submitting = false;
    }
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div class="overlay" onclick={onCancel} role="presentation">
  <!-- svelte-ignore a11y_interactive_supports_focus a11y_click_events_have_key_events -->
  <div class="modal" onclick={(e) => e.stopPropagation()} role="dialog">
    <h2>New Project</h2>

    <form onsubmit={handleSubmit}>
      <div class="field">
        <label for="projectName">Project Name</label>
        <input
          id="projectName"
          type="text"
          bind:value={name}
          placeholder="My Project"
          required
        />
      </div>

      {#if error}
        <p class="error">{error}</p>
      {/if}

      <div class="actions">
        <button type="button" class="cancel" onclick={onCancel}>Cancel</button>
        <button type="submit" disabled={submitting}>
          {submitting ? 'Creating...' : 'Create Project'}
        </button>
      </div>
    </form>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 200;
  }

  .modal {
    background: white;
    border-radius: 8px;
    padding: 24px;
    width: 380px;
    max-width: 90vw;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }

  h2 {
    margin: 0 0 16px;
    font-size: 1.25rem;
    color: #333;
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
    margin: 0 0 12px;
  }

  .actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }

  button {
    padding: 8px 16px;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
  }

  button[type="submit"] {
    background: #4a90d9;
    color: white;
    border: none;
  }

  button[type="submit"]:hover:not(:disabled) {
    background: #357abd;
  }

  button[type="submit"]:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .cancel {
    background: none;
    border: 1px solid #ccc;
    color: #555;
  }

  .cancel:hover {
    background: #f5f5f5;
  }
</style>
