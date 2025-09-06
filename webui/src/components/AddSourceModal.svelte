<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";

  interface Props {
    onCloseRequest: () => void;
    onSubmit?: (data: CreateSourceData) => void;
  }

  interface CreateSourceData {
    name: string;
    display_name: string;
    parameter: string;
    countback: number;
    is_disabled: boolean;
    schedules: string[];
  }

  const { onCloseRequest, onSubmit }: Props = $props();

  let formData = $state<CreateSourceData>({
    name: "",
    display_name: "",
    parameter: "",
    countback: 50,
    is_disabled: false,
    schedules: [],
  });

  let scheduleInput = $state("");
  let isSubmitting = $state(false);

  const listSourcesResult = useQuery(["sources", "listDropdown"], async () => {
    const { createClient } = await import("@connectrpc/connect");
    const { createConnectTransport } = await import("@connectrpc/connect-web");
    const { SourceService } = await import("../gen/claw/v1/source_service_pb");
    const transport = createConnectTransport({
      baseUrl: import.meta.env.DEV
        ? "http://localhost:8000"
        : import.meta.env.BASE_URL,
    });
    const client = createClient(SourceService, transport);
    return client.listAvailableSources({});
  });

  function handleSubmit() {
    if (
      !formData.name.trim() ||
      !formData.display_name.trim() ||
      !formData.parameter.trim()
    ) {
      return;
    }

    isSubmitting = true;
    onSubmit?.(formData);
    isSubmitting = false;
  }

  function addSchedule() {
    if (scheduleInput.trim()) {
      formData.schedules = [...formData.schedules, scheduleInput.trim()];
      scheduleInput = "";
    }
  }

  function removeSchedule(index: number) {
    formData.schedules = formData.schedules.filter((_, i) => i !== index);
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      onCloseRequest();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="modal-overlay" role="button" onclick={onCloseRequest}>
  <div class="modal-content" onclick={(e) => e.stopPropagation()}>
    <div class="modal-header">
      <h2>Add New Source</h2>
      <button class="close-btn" onclick={onCloseRequest}>×</button>
    </div>

    <form
      class="modal-form"
      onsubmit={(e) => {
        e.preventDefault();
        handleSubmit();
      }}
    >
      <div class="form-group">
        <label for="name">Source <span class="required">*</span></label>
        {#if $listSourcesResult.isLoading}
          <select id="name" disabled>
            <option>Loading sources...</option>
          </select>
        {:else if $listSourcesResult.isError}
          <select id="name" disabled>
            <option>Error loading sources</option>
          </select>
        {:else}
          <select id="name" bind:value={formData.name} required>
            <option value="" disabled selected>Select a source</option>
            {#each $listSourcesResult.data?.sources || [] as source}
              <option value={source.name}>{source.displayName}</option>
            {/each}
          </select>
        {/if}
        <input
          id="name"
          type="text"
          bind:value={formData.name}
          placeholder="e.g., reddit, booru"
          required
        />
        <small>Internal source identifier</small>
      </div>

      <div class="form-group">
        <label for="display_name"
          >Display Name <span class="required">*</span></label
        >
        <input
          id="display_name"
          type="text"
          bind:value={formData.display_name}
          placeholder="e.g., Reddit Wallpapers"
          required
        />
        <small>Human-readable name shown in UI</small>
      </div>

      <div class="form-group">
        <label for="parameter">Parameters <span class="required">*</span></label
        >
        <textarea
          id="parameter"
          bind:value={formData.parameter}
          placeholder="Configuration parameters (JSON, comma-separated values, etc.)"
          rows="3"
          required
        ></textarea>
        <small>Source-specific configuration parameters</small>
      </div>

      <div class="form-group">
        <label for="countback">Count Back</label>
        <input
          id="countback"
          type="number"
          bind:value={formData.countback}
          min="0"
          step="1"
        />
        <small
          >Number of posts to look back when searching (0 = unlimited)</small
        >
      </div>

      <div class="form-group">
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={formData.is_disabled} />
          Start disabled
        </label>
        <small>Source will be created but won't run automatically</small>
      </div>

      <div class="form-group">
        <label for="schedule-input">Schedules</label>
        <div class="schedule-input-group">
          <input
            id="schedule-input"
            type="text"
            bind:value={scheduleInput}
            placeholder="0 */6 * * * (cron expression)"
            onkeydown={(e) =>
              e.key === "Enter" && (e.preventDefault(), addSchedule())}
          />
          <button
            type="button"
            onclick={addSchedule}
            disabled={!scheduleInput.trim()}
          >
            Add
          </button>
        </div>
        <small>Cron expressions for automated runs</small>

        {#if formData.schedules.length > 0}
          <div class="schedule-list">
            {#each formData.schedules as schedule, index}
              <div class="schedule-item">
                <code>{schedule}</code>
                <button type="button" onclick={() => removeSchedule(index)}
                  >×</button
                >
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <div class="form-actions">
        <button type="button" class="btn-secondary" onclick={onCloseRequest}>
          Cancel
        </button>
        <button
          type="submit"
          class="btn-primary"
          disabled={isSubmitting ||
            !formData.name.trim() ||
            !formData.display_name.trim() ||
            !formData.parameter.trim()}
        >
          {isSubmitting ? "Creating..." : "Create Source"}
        </button>
      </div>
    </form>
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: hsla(0, 0%, 0%, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
  }

  .modal-content {
    background-color: hsl(0, 0%, 18%);
    border-radius: 8px;
    max-width: 600px;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 25px 50px hsla(0, 0%, 0%, 0.5);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 2rem;
    border-bottom: 1px solid hsl(0, 0%, 25%);
  }

  .modal-header h2 {
    margin: 0;
    color: hsl(0, 0%, 100%);
    font-size: 1.5rem;
    font-weight: 600;
  }

  .close-btn {
    background: none;
    border: none;
    color: hsl(0, 0%, 67%);
    font-size: 2rem;
    cursor: pointer;
    padding: 0;
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .close-btn:hover {
    color: hsl(0, 0%, 100%);
    background-color: hsl(0, 0%, 25%);
  }

  .modal-form {
    padding: 2rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group:last-of-type {
    margin-bottom: 2rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: hsl(0, 0%, 100%);
    font-weight: 500;
    font-size: 0.9rem;
  }

  .required {
    color: hsl(0, 100%, 70%);
  }

  input,
  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid hsl(0, 0%, 37%);
    border-radius: 6px;
    background-color: hsl(0, 0%, 24%);
    color: hsl(0, 0%, 100%);
    font-size: 0.9rem;
    transition: border-color 0.2s ease;
  }

  input:focus,
  textarea:focus {
    outline: none;
    border-color: hsl(235, 100%, 65%);
  }

  textarea {
    resize: vertical;
    min-height: 3rem;
  }

  small {
    display: block;
    margin-top: 0.25rem;
    color: hsl(0, 0%, 67%);
    font-size: 0.8rem;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    width: auto;
    margin: 0;
  }

  .schedule-input-group {
    display: flex;
    gap: 0.5rem;
  }

  .schedule-input-group input {
    flex: 1;
  }

  .schedule-input-group button {
    padding: 0.75rem 1rem;
    border: none;
    border-radius: 6px;
    background-color: hsl(235, 100%, 65%);
    color: hsl(0, 0%, 100%);
    font-size: 0.9rem;
    cursor: pointer;
    transition: background-color 0.2s ease;
    white-space: nowrap;
  }

  .schedule-input-group button:hover:not(:disabled) {
    background-color: hsl(235, 85%, 60%);
  }

  .schedule-input-group button:disabled {
    background-color: hsl(0, 0%, 37%);
    cursor: not-allowed;
  }

  .schedule-list {
    margin-top: 0.75rem;
  }

  .schedule-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.5rem 0.75rem;
    background-color: hsl(0, 0%, 24%);
    border-radius: 4px;
    margin-bottom: 0.5rem;
  }

  .schedule-item code {
    color: hsl(160, 100%, 80%);
    font-family: "Monaco", "Menlo", "Ubuntu Mono", monospace;
    font-size: 0.85rem;
  }

  .schedule-item button {
    background: none;
    border: none;
    color: hsl(0, 0%, 67%);
    cursor: pointer;
    font-size: 1.2rem;
    padding: 0.25rem;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .schedule-item button:hover {
    color: hsl(0, 100%, 70%);
    background-color: hsl(0, 0%, 30%);
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    padding-top: 1rem;
    border-top: 1px solid hsl(0, 0%, 25%);
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.75rem 1.5rem;
    border-radius: 6px;
    border: none;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-primary {
    background-color: hsl(235, 100%, 65%);
    color: hsl(0, 0%, 100%);
  }

  .btn-primary:hover:not(:disabled) {
    background-color: hsl(235, 85%, 60%);
  }

  .btn-primary:disabled {
    background-color: hsl(0, 0%, 37%);
    cursor: not-allowed;
  }

  .btn-secondary {
    background-color: hsl(0, 0%, 24%);
    color: hsl(0, 0%, 67%);
    border: 1px solid hsl(0, 0%, 37%);
  }

  .btn-secondary:hover {
    background-color: hsl(0, 0%, 30%);
    color: hsl(0, 0%, 100%);
  }
</style>

