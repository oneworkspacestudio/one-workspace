import { Fragment } from 'react'
import { Alert, Badge, Button, Card, Grid, Heading, Image, Input, Link, Page, Row, Spacer, Text } from "../main"

export function Console({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Row>
        <Heading level={1}>Web Search</Heading>
        <Spacer />
        <Link href={data.workspace} variant="apps">Apps</Link>
      </Row>
      <Text>Private browser research for your LLMs</Text>
      {data.signedOut ? (
        <Card>
          <Heading level={2}>Connect your account</Heading>
          <Text>Use your Workspace owner account to manage Web Search.</Text>
          <Link href={data.login}>Continue with MyAccounts</Link>
        </Card>
      ) : (
        <>
          <Card>
            <Row>
              <Heading level={2}>Browser service</Heading>
              <Badge>{data.state}</Badge>
            </Row>
            <Text>Fara 1.5 4B · Local model · Private browser profile</Text>
            <Text>{data.queueLabel}</Text>
            <Text>{data.setupLabel}</Text>
            <Button onClick={() => onAction("install")} disabled={data.installDisabled}>Install local model (3.4 GB)</Button>
            <Row>
              <Button onClick={() => onAction("enable")} disabled={data.enableDisabled}>Turn on</Button>
              <Button onClick={() => onAction("disable")} disabled={data.disabled}>Turn off</Button>
            </Row>
            <Alert>{data.notice}</Alert>
          </Card>
          <Grid>
            <Card>
              <Heading level={2}>Ask a question</Heading>
              <Input name="query" label="Research question" type="textarea" placeholder="What would you like to find on the web?" />
              <Button onClick={() => onAction("search")} disabled={data.notReady}>Queue research</Button>
              {data.jobs.map((item, i) => (
                <Card key={i}>
                  <Row>
                    <Text>{item.query}</Text>
                    <Badge>{item.status}</Badge>
                  </Row>
                  <Text>{item.error}</Text>
                  <Text>{item.answer}</Text>
                  {item.sources.map((item, i) => (
                    <Fragment key={i}>
                      <Link href={item.url}>{item.title}</Link>
                      <Text>{item.excerpt}</Text>
                    </Fragment>
                  ))}
                  <Button onClick={() => onAction("cancel")} item={item.id}>Cancel request</Button>
                </Card>
              ))}
            </Card>
            <Card>
              <Heading level={2}>API access</Heading>
              <Text>Credentials grant trusted server apps access to this browser profile. Keep them on your server.</Text>
              <Input name="name" label="Credential name" placeholder="My research agent" />
              <Button onClick={() => onAction("create")}>Create credential</Button>
              {data.token && (
                <>
                  <Alert>Copy this credential now. It is shown once.</Alert>
                  <Text variant="pre">{data.token}</Text>
                  <Button onClick={() => onAction("dismiss")}>I have saved it</Button>
                </>
              )}
              <Text variant="pre">{data.apiExample}</Text>
              {data.grants.map((item, i) => (
                <Row key={i}>
                  <Text>{item.name}</Text>
                  <Button onClick={() => onAction("revoke")} item={item.hash}>Revoke</Button>
                </Row>
              ))}
            </Card>
          </Grid>
          <Card id="remote-browser">
            <Heading level={2}>Remote browser</Heading>
            <Text>Resolve sign-in prompts and browser issues in the same private profile. Taking control pauses research. No local model is needed.</Text>
            <Text>{data.remoteLabel}</Text>
            <Row>
              <Button onClick={() => onAction("openRemote")} disabled={data.remoteOpenDisabled}>Open remote browser</Button>
              <Button onClick={() => onAction("reconnectRemote")} disabled={data.remoteClosed}>Reconnect</Button>
              <Button onClick={() => onAction("closeRemote")} disabled={data.remoteClosed}>Release control</Button>
            </Row>
            {data.remoteOpen && (
              <>
                <Text>Click or tap the screen to focus it, then type, paste, drag or scroll. Interrupted research must be submitted again. Closing this view releases the browser back to research.</Text>
                <Input name="address" label="Browser address" placeholder="https://example.com" />
                <Row>
                  <Button onClick={() => onAction("navigate")} disabled={data.remoteNotReady}>Go</Button>
                  <Button onClick={() => onAction("back")} disabled={data.remoteNotReady}>Back</Button>
                  <Button onClick={() => onAction("forward")} disabled={data.remoteNotReady}>Forward</Button>
                  <Button onClick={() => onAction("reload")} disabled={data.remoteNotReady}>Reload</Button>
                  <Button onClick={() => onAction("newTab")} disabled={data.remoteNotReady}>New tab</Button>
                  <Button onClick={() => onAction("closeTab")} disabled={data.remoteNotReady}>Close tab</Button>
                  <Button onClick={() => onAction("fullscreen")} disabled={data.remoteNotReady}>Full screen</Button>
                </Row>
                {data.tabs.map((item, i) => (
                  <Button key={i} onClick={() => onAction("selectTab")} item={item.id} disabled={item.active}>{item.label}</Button>
                ))}
                <Text>{data.remotePage}</Text>
                {data.remoteImage && (
                  <Image variant="remote" src={data.remoteImage} alt="Remote browser screen" />
                )}
                <Input name="remoteText" label="Private text to type" type="password" placeholder="Type or paste text, then choose Send text" />
                <Row>
                  <Button onClick={() => onAction("sendText")} disabled={data.remoteNotReady}>Send text</Button>
                  <Button onClick={() => onAction("enter")} disabled={data.remoteNotReady}>Enter</Button>
                  <Button onClick={() => onAction("tab")} disabled={data.remoteNotReady}>Tab</Button>
                  <Button onClick={() => onAction("scrollUp")} disabled={data.remoteNotReady}>Scroll up</Button>
                  <Button onClick={() => onAction("scrollDown")} disabled={data.remoteNotReady}>Scroll down</Button>
                </Row>
                <Text>This view controls the research browser, not your PC desktop. Clipboard and operating-system shortcuts stay blocked. Release control when finished; disconnected sessions expire after 90 seconds.</Text>
              </>
            )}
          </Card>
          <Text>Requests and results expire after 24 hours. Sign-in challenges pause research.</Text>
        </>
      )}
    </Page>
  )
}
