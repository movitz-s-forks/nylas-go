package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestThreads(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"any_email":           {"a@example.com,b@example.com"},
		"bcc":                 {"f@example.com"},
		"cc":                  {"e@example.com"},
		"filename":            {"filename"},
		"from":                {"d@example.com"},
		"in":                  {"in"},
		"last_message_after":  {"4"},
		"last_message_before": {"3"},
		"limit":               {"1"},
		"offset":              {"2"},
		"started_after":       {"6"},
		"started_before":      {"5"},
		"subject":             {"subject"},
		"to":                  {"c@example.com"},
		"unread":              {"true"},
		"view":                {"ids"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(threadsJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Threads(context.Background(), &ThreadsOptions{
		AnyEmail:          []string{"a@example.com", "b@example.com"},
		BCC:               "f@example.com",
		CC:                "e@example.com",
		Filename:          "filename",
		From:              "d@example.com",
		In:                "in",
		LastMessageAfter:  4,
		LastMessageBefore: 3,
		Limit:             1,
		Offset:            2,
		StartedAfter:      6,
		StartedBefore:     5,
		Subject:           "subject",
		To:                "c@example.com",
		Unread:            true,
		View:              "ids",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Thread{
		{
			AccountID:             "1234***",
			DraftIDs:              []string{},
			FirstMessageTimestamp: 1557950729,
			Folders: []Folder{
				{
					DisplayName: "Inbox",
					ID:          "4567****",
					Name:        "inbox",
				},
			},
			HasAttachments:               false,
			ID:                           "4312****",
			LastMessageReceivedTimestamp: 1557950729,
			LastMessageSentTimestamp:     0,
			LastMessageTimestamp:         1557950729,
			MessageIDs: []string{
				"5634***",
			},
			Object: "thread",
			Participants: []Participant{
				{
					Email: "no-reply@cc.yahoo-inc.com",
					Name:  "Yahoo",
				},
				{
					Email: "james****@yahoo.com",
					Name:  "",
				},
			},
			Snippet: "Hi James, welcome.",
			Starred: false,
			Subject: "Security settings changed on your Yahoo account",
			Unread:  false,
			Version: 1,
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Threads: (-got +want):\n%s", diff)
	}
}

func TestThread(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{}
	id := "8r5awu0esbg8ct3wg5rj5sifp"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertQueryParams(t, r, wantQuery)
		if r.URL.Path != "/threads/"+id {
			t.Errorf("unexpected path: %v", r.URL.Path)
		}
		_, _ = w.Write(getThreadJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Thread(context.Background(), id, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Thread{
		AccountID:             "crkr5ct7aa3edvipotbj2****",
		DraftIDs:              []string{},
		FirstMessageTimestamp: 1579611155,
		HasAttachments:        false,
		ID:                    "8r5awu0esbg8ct3wg5rj5sifp",
		Labels: []Label{
			{
				DisplayName: "Important",
				ID:          "a1ytpbvawxfaqua671478g1q0",
				Name:        "important",
			},
			{
				DisplayName: "Inbox",
				ID:          "atamsqdb355jqyj0zhhatu3ao",
				Name:        "inbox",
			},
		},
		LastMessageReceivedTimestamp: 1579611155,
		LastMessageSentTimestamp:     0,
		LastMessageTimestamp:         1579611155,
		MessageIDs: []string{
			"br57kcekhf1hsjq04y8aonkit",
		},
		Object: "thread",
		Participants: []Participant{
			{
				Email: "from@example.org",
				Name:  "From Name",
			},
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		Snippet: "Body",
		Starred: true,
		Subject: "Subject",
		Unread:  true,
		Version: 2,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Thread: (-got +want):\n%s", diff)
	}
}

var threadsJSON = []byte(`[
    {
	"account_id": "1234***",
	"draft_ids": [],
	"first_message_timestamp": 1557950729,
	"folders": [
	    {
		"display_name": "Inbox",
		"id": "4567****",
		"name": "inbox"
	    }
	],
	"has_attachments": false,
	"id": "4312****",
	"last_message_received_timestamp": 1557950729,
	"last_message_sent_timestamp": null,
	"last_message_timestamp": 1557950729,
	"message_ids": [
	    "5634***"
	],
	"object": "thread",
	"participants": [
	    {
		"email": "no-reply@cc.yahoo-inc.com",
		"name": "Yahoo"
	    },
	    {
		"email": "james****@yahoo.com",
		"name": ""
	    }
	],
	"snippet": "Hi James, welcome.",
	"starred": false,
	"subject": "Security settings changed on your Yahoo account",
	"unread": false,
	"version": 1
    }
]`)

var getThreadJSON = []byte(`{
  "account_id": "crkr5ct7aa3edvipotbj2****",
  "draft_ids": [],
  "first_message_timestamp": 1579611155,
  "has_attachments": false,
  "id": "8r5awu0esbg8ct3wg5rj5sifp",
  "labels": [
    {
      "display_name": "Important",
      "id": "a1ytpbvawxfaqua671478g1q0",
      "name": "important"
    },
    {
      "display_name": "Inbox",
      "id": "atamsqdb355jqyj0zhhatu3ao",
      "name": "inbox"
    }
  ],
  "last_message_received_timestamp": 1579611155,
  "last_message_sent_timestamp": null,
  "last_message_timestamp": 1579611155,
  "message_ids": [
    "br57kcekhf1hsjq04y8aonkit"
  ],
  "object": "thread",
  "participants": [
    {
      "email": "from@example.org",
      "name": "From Name"
    },
    {
      "email": "to@example.org",
      "name": "To Name"
    }
  ],
  "snippet": "Body",
  "starred": true,
  "subject": "Subject",
  "unread": true,
  "version": 2
}`)
