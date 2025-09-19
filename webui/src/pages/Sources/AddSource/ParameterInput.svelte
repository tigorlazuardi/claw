<script lang="ts">
  import type {
    AvailableSource,
    ValidateSourceParametersResponse,
  } from "#/gen/claw/v1/source_service_pb";
  import { Popover } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import { getSourceServiceClient } from "#/connectrpc";
  import IconCheck from "@lucide/svelte/icons/check";
  import { theme } from "#/store/theme";
  import { resource, watch } from "runed";

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

  const validateParameter = resource(
    () => ({ source, value }),
    async (val, prev, { data, signal }) => {
      if (!val.value.trim()) return;
      if (val.value.trim() === prev?.value.trim())
        return data as ValidateSourceParametersResponse;
      return getSourceServiceClient({ signal }).then((client) =>
        client.validateSourceParameters({
          sourceName: val.source.name,
          parameter: val.value,
        }),
      );
    },
    {
      lazy: true,
      debounce: 500,
    },
  );

  let error = $state("");

  watch(
    () => validateParameter.current,
    (data) => {
      if (data) {
        error = "";
        if (data.transformedParameter) {
          value = data.transformedParameter;
        }
        valid = true;
      }
    },
  );
  watch(
    () => validateParameter.error,
    (err) => {
      if (err) {
        error = err.message;
        valid = false;
        return;
      }
    },
  );

  function handleValidateParamaterOnBlur() {
    if (value.trim().length === 0) {
      if (source.requireParameter) {
        error = "This field is required";
        valid = false;
        return;
      }
      error = "";
      valid = true;
      return;
    }
    validateParameter.refetch();
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
          <Popover.Content class="z-[9999] bg-transparent" data-theme={$theme}>
            <div class="card bg-base-300 border-base-200 border">
              <div class="card-body">
                <div class="card-title">Parameter Help</div>
                <div
                  class="p-4 border border-base-100"
                  style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
                >
                  {#if showParameterHelp}
                    {#await import ("#/components/MarkdownText.svelte") then { default: MarkdownText }}
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
    bind:value
    onblur={handleValidateParamaterOnBlur}
    oninput={() => {
      if (source.requireParameter) {
        valid = false;
      }
    }}
    required={source.requireParameter}
  ></textarea>
  {#if validateParameter.loading}
    <div class="alert alert-warning alert-soft">
      <div class="loading loading-spinner"></div>
      <span>Validating...</span>
    </div>
  {:else if error}
    {#await import ("#/components/MarkdownText.svelte") then { default: MarkdownText }}
      <div role="alert" class="alert alert-error alert-soft">
        <MarkdownText text={error} />
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
