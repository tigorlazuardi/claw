<script lang="ts">
  import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
  import { form as formValidator, field } from "svelte-forms";
  import { required } from "svelte-forms/validators";
  import type { M } from "../types";
  import type { CreateSourceRequest } from "../gen/claw/v1/source_service_pb";

  interface Props {
    /**
     * Callback when the modal is requesting to be closed or removed from the DOM.
     * This will be called when the user clicks the close button, cancel button,
     * or the server responds with success after form submission.
     *
     * This function will not be called if the user clicks outside the modal or press the escape key
     * to prevent accidental closure.
     */
    onCloseRequest: () => void;
  }

  const name = field("name", "", [required()]);
  const parameter = field("parameter", "", [required()]);
  const display_name = field("display_name", "", [required()]);
  const countback = field("countback", 200);
  const is_disabled = field("is_disabled", false);

  const schedules = $state<string[]>([]);
  let scheduleExpressionError = $state<string | undefined>();
  const scheduleInput = field("schedule", "", [
    async (val: string) => {
      if (val.trim() === "") {
        scheduleExpressionError = undefined;
        return { valid: true, name: "cron_format" };
      }
      const { CronExpressionParser } = await import("cron-parser");
      try {
        CronExpressionParser.parse(val);
        scheduleExpressionError = undefined;
        return { valid: true, name: "cron_format" };
      } catch (e) {
        scheduleExpressionError =
          "invalid cron expression: " + (e as Error).message;
        return { valid: false, name: "cron_format" };
      }
    },
  ]);
  const addSourceForm = formValidator(
    name,
    parameter,
    display_name,
    countback,
    is_disabled,
    scheduleInput,
  );

  const { onCloseRequest }: Props = $props();

  async function createConnectClient() {
    const { createClient } = await import("@connectrpc/connect");
    const { createConnectTransport } = await import("@connectrpc/connect-web");
    const { SourceService } = await import("../gen/claw/v1/source_service_pb");
    const transport = createConnectTransport({
      baseUrl: import.meta.env.BASE_URL,
    });
    return createClient(SourceService, transport);
  }

  async function postCreateRequest() {
    if (!$addSourceForm.valid) {
      return;
    }
    const data: M<CreateSourceRequest> = {
      name: $name.value,
      displayName: $display_name.value,
      parameter: $parameter.value,
      countback: Number($countback.value),
      isDisabled: Boolean($is_disabled.value),
      schedules: schedules,
    };
    return createConnectClient().then((client) => client.createSource(data));
  }

  let serverErrorResponse = $state<Error | null>(null);
  function handleOnSubmit() {
    return postCreateRequest()
      .then(async (res) => {
        if (!res) return;
        const queryClient = useQueryClient();
        queryClient.invalidateQueries(["sources", "list"]);
        onCloseRequest();
      })
      .catch((err: Error) => (serverErrorResponse = err));
  }

  let isSubmitting = $state(false);
  let showParameterHelp = $state(false);

  async function listSources() {
    const { createClient } = await import("@connectrpc/connect");
    const { createConnectTransport } = await import("@connectrpc/connect-web");
    const { SourceService } = await import("../gen/claw/v1/source_service_pb");
    const transport = createConnectTransport({
      baseUrl: import.meta.env.BASE_URL,
    });
    const client = createClient(SourceService, transport);
    return client.listAvailableSources({});
  }

  const listSourcesResult = useQuery(["sources", "listDropdown"], listSources);

  // Get the selected source object
  let selectedSource = $derived(
    $listSourcesResult.data?.sources?.find((s) => s.name === $name.value),
  );

  // Get parameter placeholder from selected source
  let parameterPlaceholder = $derived(
    selectedSource?.parameterPlaceholder ||
      "Configuration parameters (JSON, comma-separated values, etc.)",
  );

  // Check if parameter help should be shown
  let hasParameterHelp = $derived(
    selectedSource?.parameterHelp && selectedSource.parameterHelp.trim() !== "",
  );

  async function addSchedule() {
    await scheduleInput.validate();
    if ($addSourceForm.hasError("schedule.cron_format")) return;
    schedules.push($scheduleInput.value);
    scheduleInput.reset();
  }

  function removeSchedule(i: number) {
    schedules.splice(i, 1);
  }
</script>

<div class="modal-overlay" role="button" tabindex="0" aria-label="Close modal">
  <div
    class="modal-content"
    role="dialog"
    aria-modal="true"
    aria-labelledby="modal-title"
    tabindex="-1"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <div class="modal-header">
      <h2 id="modal-title">Add New Source</h2>
      <button class="close-btn" onclick={onCloseRequest}>×</button>
    </div>

    <form
      class="modal-form"
      onsubmit={(e) => {
        e.preventDefault();
        handleOnSubmit();
      }}
    >
      <div class="form-group">
        <label for="name">
          Source <span class="required">*</span>
        </label>
        {#if $listSourcesResult.isLoading}
          <select name={$name.name} id="name" disabled>
            <option>Loading sources...</option>
          </select>
        {:else if $listSourcesResult.isError}
          <select name={$name.name} id="name" disabled>
            <option>Error loading sources</option>
          </select>
        {:else}
          <select name={$name.name} id="name" bind:value={$name.value}>
            <option value="" disabled selected>Select a source</option>
            {#each $listSourcesResult.data?.sources || [] as source}
              <option value={source.name}>
                {source.displayName} ({source.name})
              </option>
            {/each}
          </select>
        {/if}
        <small>Choose supported source</small>
      </div>

      {#if sourceSelected}
        <div class="form-group">
          <div class="label-with-help">
            <label for="parameter">
              Parameters <span class="required">*</span>
            </label>
            {#if hasParameterHelp}
              <button
                type="button"
                class="help-btn"
                onclick={() => (showParameterHelp = !showParameterHelp)}
                aria-label="Show parameter help"
                title="Show parameter help"
              >
                {@render infoIcon()}
              </button>
            {/if}
          </div>
          <textarea
            name={$parameter.name}
            id="parameter"
            bind:value={$parameter.value}
            placeholder={parameterPlaceholder}
            rows="3"
          ></textarea>
          {#if showParameterHelp && selectedSource?.parameterHelp}
            {#await import("./MarkdownText.svelte") then { default: MarkdownText }}
              <div class="parameter-help">
                <MarkdownText text={selectedSource.parameterHelp} />
              </div>
            {/await}
          {/if}
          {@render labelOrError(
            "Source-specific configuration parameters",
            $name.errors.join(" "),
          )}
        </div>

        <div class="form-group">
          <label for={$display_name.name}>
            Display Name <span class="required">*</span>
          </label>
          <input
            id="display_name"
            name={$display_name.name}
            type="text"
            bind:value={$display_name.value}
            placeholder="e.g., Reddit Wallpapers"
          />
          <small>Human-readable name shown in UI</small>
        </div>

        <div class="form-group">
          <label for={$countback.name}>Count Back</label>
          <input
            id="countback"
            name={$countback.name}
            type="number"
            bind:value={$countback.value}
            min="0"
            step="1"
          />
          <small>
            Number of posts to look back when searching (0 = source default)
          </small>
        </div>

        <div class="form-group">
          <label class="checkbox-label" for={$is_disabled.name}>
            <input type="checkbox" bind:checked={$is_disabled.value} />
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
              bind:value={$scheduleInput.value}
              placeholder="0 */6 * * * (cron expression)"
              onkeydown={(e) =>
                e.key === "Enter" && (e.preventDefault(), addSchedule())}
            />
            <button
              type="button"
              onclick={addSchedule}
              disabled={!$scheduleInput.valid}
            >
              Add
            </button>
          </div>
          {@render labelOrError(
            "Add cron expressions to schedule automated runs",
            scheduleExpressionError,
          )}

          {#if schedules.length > 0}
            <div class="schedule-list">
              {#each schedules as schedule, index (index)}
                <div class="schedule-item">
                  <code>{schedule}</code>
                  <button type="button" onclick={() => removeSchedule(index)}>
                    ×
                  </button>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}

      <div class="form-actions">
        <button type="button" class="btn-secondary" onclick={onCloseRequest}>
          Cancel
        </button>
        <button
          type="submit"
          class="btn-primary"
          disabled={!$addSourceForm.valid}
        >
          {isSubmitting ? "Creating..." : "Create Source"}
        </button>
      </div>
    </form>
  </div>
</div>

{#snippet labelOrError(helpText: string, errString?: string)}
  {#if errString}
    <small class="error-label">{errString}</small>
  {:else}
    <small>{helpText}</small>
  {/if}
{/snippet}

{#snippet infoIcon()}
  <svg
    class="help-icon"
    fill="currentColor"
    version="1.1"
    id="Capa_1"
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 488.484 488.484"
    xml:space="preserve"
  >
    <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
    <g
      id="SVGRepo_tracerCarrier"
      stroke-linecap="round"
      stroke-linejoin="round"
    ></g>
    <g id="SVGRepo_iconCarrier">
      <g>
        <g>
          <path
            d="M244.236,0.002C109.562,0.002,0,109.565,0,244.238c0,134.679,109.563,244.244,244.236,244.244 c134.684,0,244.249-109.564,244.249-244.244C488.484,109.566,378.92,0.002,244.236,0.002z M244.236,413.619 c-93.4,0-169.38-75.979-169.38-169.379c0-93.396,75.979-169.375,169.38-169.375s169.391,75.979,169.391,169.375 C413.627,337.641,337.637,413.619,244.236,413.619z"
          ></path>
          <path
            d="M244.236,206.816c-14.757,0-26.619,11.962-26.619,26.73v118.709c0,14.769,11.862,26.735,26.619,26.735 c14.769,0,26.62-11.967,26.62-26.735V233.546C270.855,218.778,259.005,206.816,244.236,206.816z"
          ></path>
          <path
            d="M244.236,107.893c-19.949,0-36.102,16.158-36.102,36.091c0,19.934,16.152,36.092,36.102,36.092 c19.929,0,36.081-16.158,36.081-36.092C280.316,124.051,264.165,107.893,244.236,107.893z"
          ></path>
        </g>
      </g>
    </g>
  </svg>
{/snippet}

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
    max-width: 60vw;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 25px 50px hsla(0, 0%, 0%, 0.5);
  }

  @media (max-width: 500px) {
    .modal-content {
      max-width: 100vw;
    }
  }

  @media (max-width: 1200px) {
    .modal-content {
      max-width: 75vw;
    }
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
  textarea,
  select {
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
  textarea:focus,
  select:focus {
    outline: none;
    border-color: hsl(235, 100%, 65%);
  }

  select {
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3E%3Cpath fill='%23aaa' d='M4.427 9.427l3.396 3.396a.25.25 0 00.354 0l3.396-3.396A.25.25 0 0011.396 9H4.604a.25.25 0 00-.177.427z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 16px;
    padding-right: 2.5rem;
  }

  select:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  select option {
    background-color: hsl(0, 0%, 24%);
    color: hsl(0, 0%, 100%);
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

  .label-with-help {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
    gap: 0.5rem;
  }

  .help-btn {
    margin: 0;
    border-radius: 50%;
    border: none;
    color: hsl(0, 0%, 100%);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    flex-shrink: 0;
    box-shadow: 0 2px 4px hsla(0, 0%, 0%, 0.2);
    padding: 0;
  }

  .help-btn:hover {
    background-color: hsl(235, 85%, 60%);
    transform: translateY(-1px);
    box-shadow: 0 4px 8px hsla(0, 0%, 0%, 0.3);
  }

  .help-icon {
    width: 0.75rem;
    height: 0.75rem;
    display: block;
  }

  .parameter-help {
    background-color: hsl(210, 25%, 95%);
    color: hsl(210, 15%, 35%);
    padding: 0.75rem;
    border-radius: 6px;
    margin-top: 0.5rem;
    border: 1px solid hsl(210, 25%, 85%);
  }
</style>
