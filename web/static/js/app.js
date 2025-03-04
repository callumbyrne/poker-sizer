document.addEventListener("alpine:init", () => {
  // Register the planningPoker component
  Alpine.data("planningPoker", function () {
    return {
      state: "voting", // 'voting', 'revealed', 'reset'
      currentIssue: "",
      selectedCard: null,
      users: {},
      votes: {},
      isAdmin: false,
      connectionStatus: "connecting",
      websocket: null,
      notificationSound: new Audio("/static/sounds/notification.mp3"),

      init() {
        // Get initial data from the DOM (set by the server template)
        try {
          this.currentIssue = this.$el.dataset.issue || "";
          this.isAdmin = this.$el.dataset.isAdmin === "true";
          this.users = JSON.parse(this.$el.dataset.users || "{}");
          this.state = this.$el.dataset.state || "voting";
        } catch (err) {
          console.error("Error parsing initial data:", err);
        }

        // Setup WebSocket event listeners
        this.$watch("connectionStatus", (value) => {
          if (value === "connected") {
            this.sendJoinNotification();
          }
        });

        // Setup event listeners for HTMX WebSocket events
        this.$el.addEventListener("ws:open", () => {
          this.connectionStatus = "connected";
        });

        this.$el.addEventListener("ws:close", () => {
          this.connectionStatus = "disconnected";
        });

        this.$el.addEventListener("ws:message", (event) => {
          const data = JSON.parse(event.detail.message);
          this.handleMessage(data);
        });

        // Auto-copy room link on creation
        if (
          this.isAdmin &&
          !localStorage.getItem(`room-${this.$el.dataset.roomId}-notified`)
        ) {
          this.copyRoomLink();
          localStorage.setItem(
            `room-${this.$el.dataset.roomId}-notified`,
            "true",
          );
        }
      },

      handleMessage(data) {
        switch (data.type) {
          case "user_joined":
            this.users[data.payload.id] = data.payload;
            this.playNotification();
            this.showToast(`${data.payload.name} joined the room`);
            break;

          case "user_left":
            delete this.users[data.payload.id];
            this.showToast(`${data.payload.name} left the room`);
            break;

          case "vote_submitted":
            this.votes[data.payload.userId] = data.payload;

            // If all users have voted and I'm the admin, show a notification
            if (this.isAdmin && this.allUsersVoted) {
              this.showToast("All users have voted!");
              this.playNotification();
            }
            break;

          case "votes_revealed":
            this.state = "revealed";
            this.votes = data.payload.votes;
            this.playNotification();
            break;

          case "room_reset":
            this.state = "voting";
            this.votes = {};
            this.selectedCard = null;
            this.showToast("Voting has been reset");
            break;

          case "issue_updated":
            this.currentIssue = data.payload.issue;
            this.showToast("Issue updated");
            break;
        }
      },

      selectCard(value) {
        // Don't reselect the same card
        if (this.selectedCard === value) return;

        this.selectedCard = value;

        // Send the vote via WebSocket
        const message = {
          type: "submit_vote",
          payload: {
            value: value,
          },
        };
        this.$dispatch("ws:send", JSON.stringify(message));

        // Add card flip animation
        const cards = document.querySelectorAll(".poker-card");
        cards.forEach((card) => {
          if (!card.classList.contains("selected")) {
            card.style.transform = "scale(0.95)";
            setTimeout(() => {
              card.style.transform = "";
            }, 300);
          }
        });
      },

      revealVotes() {
        // Animation preparation
        document.querySelectorAll(".poker-card").forEach((card) => {
          card.classList.add("disabled");
        });

        const message = {
          type: "reveal_votes",
          payload: {},
        };
        this.$dispatch("ws:send", JSON.stringify(message));
      },

      resetVoting() {
        const message = {
          type: "reset_voting",
          payload: {},
        };
        this.$dispatch("ws:send", JSON.stringify(message));
      },

      updateIssue() {
        if (!this.currentIssue.trim()) return;

        const message = {
          type: "update_issue",
          payload: {
            issue: this.currentIssue,
          },
        };
        this.$dispatch("ws:send", JSON.stringify(message));
      },

      sendJoinNotification() {
        const message = {
          type: "user_joined_notification",
          payload: {},
        };
        this.$dispatch("ws:send", JSON.stringify(message));
      },

      copyRoomLink() {
        const roomLink = window.location.href;
        navigator.clipboard
          .writeText(roomLink)
          .then(() => {
            this.showToast("Room link copied to clipboard");
          })
          .catch((err) => {
            console.error("Could not copy room link:", err);
          });
      },

      showToast(message) {
        // Create toast element if it doesn't exist
        let toast = document.getElementById("poker-toast");
        if (!toast) {
          toast = document.createElement("div");
          toast.id = "poker-toast";
          toast.className = "toast";
          document.body.appendChild(toast);
        }

        // Set message and show
        toast.textContent = message;
        toast.classList.add("show");

        // Hide after 3 seconds
        setTimeout(() => {
          toast.classList.remove("show");
        }, 3000);
      },

      playNotification() {
        // Only play if sound exists and we're not in a meeting
        if (this.notificationSound && !document.hidden) {
          this.notificationSound.play().catch((err) => {
            // Ignore autoplay errors
            console.log("Could not play notification sound");
          });
        }
      },

      get results() {
        if (this.state !== "revealed") return [];

        const voteCounts = {};
        Object.values(this.votes).forEach((vote) => {
          if (!voteCounts[vote.value]) {
            voteCounts[vote.value] = 0;
          }
          voteCounts[vote.value]++;
        });

        return Object.entries(voteCounts)
          .map(([value, count]) => ({
            value,
            count,
          }))
          .sort((a, b) => {
            // Sort numeric values numerically, and '?' at the end
            if (a.value === "?") return 1;
            if (b.value === "?") return -1;
            return Number(a.value) - Number(b.value);
          });
      },

      get average() {
        if (this.state !== "revealed") return null;

        const numericVotes = Object.values(this.votes)
          .map((vote) => {
            // Ensure only numeric values are used
            return isNaN(Number(vote.value)) ? null : Number(vote.value);
          })
          .filter((value) => value !== null);

        if (numericVotes.length === 0) return null;

        return (
          numericVotes.reduce((sum, value) => sum + value, 0) /
          numericVotes.length
        );
      },

      get allUsersVoted() {
        const userCount = Object.keys(this.users).length;
        const voteCount = Object.keys(this.votes).length;
        return userCount > 0 && voteCount === userCount;
      },

      get votingProgress() {
        const userCount = Object.keys(this.users).length;
        const voteCount = Object.keys(this.votes).length;
        return userCount > 0 ? Math.round((voteCount / userCount) * 100) : 0;
      },
    };
  });
});

// Additional CSS for toasts
document.addEventListener("DOMContentLoaded", () => {
  const style = document.createElement("style");
  style.textContent = `
        .toast {
            visibility: hidden;
            min-width: 250px;
            background-color: #333;
            color: #fff;
            text-align: center;
            border-radius: 8px;
            padding: 16px;
            position: fixed;
            z-index: 1000;
            left: 50%;
            bottom: 30px;
            transform: translateX(-50%);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
            opacity: 0;
            transition: visibility 0s 0.5s, opacity 0.5s ease;
        }
        
        .toast.show {
            visibility: visible;
            opacity: 1;
            transition: opacity 0.5s ease;
        }
    `;
  document.head.appendChild(style);
});
