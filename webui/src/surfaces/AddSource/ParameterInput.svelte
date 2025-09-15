<script lang="ts">
  import type {
    AvailableSource,
    ValidateSourceParametersRequest,
  } from "../../gen/claw/v1/source_service_pb";
  import { Popover } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import { useQuery, useQueryClient } from "@sveltestack/svelte-query";
  import { getSourceServiceClient } from "../../connectrpc";
  import type { M } from "../../types";
  import IconCheck from "@lucide/svelte/icons/check";

  const queryClient = useQueryClient();

  interface Props {
    source: AvailableSource;
    value?: string;
    valid?: boolean;
  }
  let {
    source,
    value = $bindable(""),
    valid = $bindable(!source.requireParameter),
  }: Props = $props();
  const hasParameterHelp = source.parameterHelp.trim().length > 0;
  let showParameterHelp = $state(false);

  let validateParameterRequest: M<ValidateSourceParametersRequest> = $derived({
    sourceName: source.name,
    parameter: value,
  });

  let parameterInputField: HTMLTextAreaElement | undefined = $state();

  const validateParameterQuery = useQuery(
    ["source", "add", "validateParameter"],
    () =>
      getSourceServiceClient().then((client) =>
        client.validateSourceParameters(validateParameterRequest),
      ),
    {
      onSuccess(data) {
        parameterInputField?.setCustomValidity("");
        if (data.transformedParameter) {
          value = data.transformedParameter;
          valid = true;
        }
      },
      onError(err: Error) {
        valid = false;
        parameterInputField?.setCustomValidity(err.message);
      },
      enabled: false, // Only run on demand
      retry: false,
    },
  );

  function handleValidateParamaterOnBlur() {
    if (valid) {
      return;
    }
    if (source.requireParameter && value.trim().length === 0) {
      parameterInputField?.setCustomValidity("This field is required.");
      valid = false;
      return;
    }
    if (value.trim().length === 0) {
      // No need to validate empty optional parameters
      parameterInputField?.setCustomValidity("");
      valid = true;
      return;
    }
    queryClient.cancelQueries(["source", "add", "validateParameter"]);
    $validateParameterQuery.refetch();
  }

  const allOk = $derived(
    valid && (!source.requireParameter || value.trim().length > 0),
  );
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span class={{ "text-success": allOk }}>Parameters</span>
    {#if source.requireParameter}
      <span class="text-error">*</span>
    {/if}
    {#if hasParameterHelp}
      <Popover.Root onOpenChange={(v) => (showParameterHelp = v)}>
        <Popover.Trigger class="btn btn-square btn-ghost btn-xs" type="button">
          <IconInfo />
        </Popover.Trigger>
        <Popover.Portal>
          <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula">
            <div class="card bg-base-300 border-base-200 border">
              <div class="card-body">
                <div class="card-title">Parameter Help</div>
                <div
                  class="p-4 border border-base-100"
                  style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
                >
                  {#if showParameterHelp}
                    {#await import ("../../components/MarkdownText.svelte") then { default: MarkdownText }}
                      <MarkdownText text={source.parameterHelp} />
                    {/await}
                  {/if}
                </div>
              </div>
            </div>
          </Popover.Content>
        </Popover.Portal>
      </Popover.Root>
    {/if}
  </legend>
  <textarea
    class={{
      "textarea h-[3rem] w-full": true,
      "textarea-success": allOk,
      "text-success": allOk,
    }}
    placeholder={source.parameterPlaceholder ||
      "Configuration parameters (JSON, comma-separated values, etc.)"}
    bind:this={parameterInputField}
    bind:value
    onblur={handleValidateParamaterOnBlur}
    oninput={() => {
      if (source.requireParameter) {
        valid = false;
      }
    }}
    disabled={$validateParameterQuery.isFetching}
    required={source.requireParameter}
  ></textarea>
  {#if $validateParameterQuery.isFetching}
    <div class="alert alert-warning alert-soft">
      <div class="loading loading-spinner"></div>
      <span>Validating...</span>
    </div>
  {:else if parameterInputField?.validity.customError}
    {#await import ("../../components/MarkdownText.svelte") then { default: MarkdownText }}
      <div role="alert" class="alert alert-error alert-soft">
        <MarkdownText text={parameterInputField?.validationMessage || ""} />
      </div>
    {/await}
  {:else if valid && source.requireParameter}
    <div class="alert alert-success alert-soft">
      <IconCheck />
      <span>Looks good!</span>
    </div>
  {:else if source.requireParameter}
    <span class="label">Required</span>
  {:else}
    <span class="label">Optional</span>
  {/if}
</fieldset>
