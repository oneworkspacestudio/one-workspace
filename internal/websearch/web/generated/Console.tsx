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
          <Card>
            <Heading level={2}>Browser preview</Heading>
            <Text>The browser keeps its own cookies and sessions. Live research controls it one request at a time.</Text>
            <Button onClick={() => onAction("preview")} disabled={data.notReady}>Refresh preview</Button>
            {data.preview && (
              <Image src={data.preview} alt="Private research browser screenshot" />
            )}
          </Card>
          <Text>Requests and results expire after 24 hours. Sign-in challenges pause research.</Text>
        </>
      )}
    </Page>
  )
}
