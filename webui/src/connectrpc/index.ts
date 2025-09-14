export type SourceServiceClient = Awaited<
  ReturnType<typeof createSourceServiceClient>
>;
let sourceServiceClient: SourceServiceClient;
async function createSourceServiceClient() {
  const { createClient } = await import("@connectrpc/connect");
  const { createConnectTransport } = await import("@connectrpc/connect-web");
  const { SourceService } = await import("../gen/claw/v1/source_service_pb");
  const transport = createConnectTransport({
    baseUrl: import.meta.env.BASE_URL,
    fetch(input, init) {
      const cancelToken = new AbortController().signal;
      const promise = fetch(input, {
        ...init,
        credentials: "include",
        signal: cancelToken,
      });
      //@ts-ignore svelte query will use the cancel method if it exists to cancel requests.
      promise.cancel = () => cancelToken.abort();
      return promise;
    },
  });
  return createClient(SourceService, transport);
}

export async function getSourceServiceClient() {
  if (!sourceServiceClient) {
    sourceServiceClient = await createSourceServiceClient();
  }
  return sourceServiceClient;
}
