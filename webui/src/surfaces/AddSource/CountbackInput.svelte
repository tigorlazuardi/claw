<script lang="ts">
  import { Popover } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import type { AvailableSource } from "../../gen/claw/v1/source_service_pb";
  interface Props {
    value: number;
    source: AvailableSource;
  }

  let { value = $bindable(0), source }: Props = $props();

  const allOk = $derived(value > 0);
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span class={{ "text-success": allOk }}>Count Back</span>
    <Popover.Root>
      <Popover.Trigger class="btn btn-square btn-ghost btn-xs" type="button">
        <IconInfo />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula">
          <div class="card bg-base-300 border-base-200 border">
            <div class="card-body">
              <div class="card-title flex-col items-start">
                <span>Count Back</span>
                <span class="text-xs">
                  Souce: {source.displayName} ({source.name})
                </span>
                <span class="text-xs">
                  Default Value: {source.defaultCountback}
                </span>
              </div>
              <div
                class="prose p-4 border border-base-100"
                style="box-shadow: inset 0 6px 12px rgba(0, 0, 0, 0.15), inset 0 2px 4px rgba(0, 0, 0, 0.1);"
              >
                <p>
                  Count Back lookup the number of "Items" to look up for and
                  sets an upper limit for Claw to stop looking for more images.
                </p>
                <p>
                  However, "Items" do not mean "Images". It can be something
                  like a post or forum entry, but if it's not an image, the post
                  will be skipped, however it still counts towards the count
                  back limit.
                </p>
                <p>
                  The main purpose for Count Back is to avoid API rate limit.
                </p>

                <p>
                  Source Default Value usually expects for the same Source Kind
                  to not have overlapping schedules and takes as many images
                  possible under the API limit.
                </p>
              </div>
            </div>
          </div>
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  </legend>
  <input
    class={{
      "input w-full": true,
      "input-success": allOk,
      "text-success": allOk,
    }}
    type="number"
    step="1"
    min="0"
    bind:value
  />
  <p
    class={{
      label: true,
      "text-success": allOk,
    }}
  >
    The number of items to look up for. If value is 0, source default will be
    used
  </p>
</fieldset>
