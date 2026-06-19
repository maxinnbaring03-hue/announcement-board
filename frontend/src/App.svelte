<script>
  import { onMount } from 'svelte';

  let currentView = 'public';

  /** @type {any[]} */
  let announcements = [];

  let username = '';
  let password = '';

  let newUsername = '';
  let newPassword = '';
  let confirmPassword = '';

  let draftTitle = '';
  let draftMessage = '';

  onMount(async () => {
    if (window.location.hash === '#create-admin') {
      currentView = 'register';
    }
    await fetchAnnouncements();
  });

  async function fetchAnnouncements() {
    try {
      const response = await fetch('http://localhost:8080/api/announcements');
      const data = await response.json();
      announcements = data || []; 
    } catch (error) {
      console.error("Connection error:", error);
    }
  }

  function goToLogin() { currentView = 'login'; }

  function cancelLogin() {
    currentView = 'public';
    username = ''; password = '';
    window.location.hash = ''; 
  }

  async function handleLogin() {
    if (!username || !password) {
      alert("Please enter both username and password.");
      return;
    }
    try {
      const res = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      if (res.ok) {
        currentView = 'admin'; 
      } else {
        alert("Invalid username or password. Please try again.");
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("Failed to reach the server.");
    }
  }

  async function handleRegister() {
    if (!newUsername || !newPassword) {
      alert("Please fill out all fields.");
      return;
    }
    if (newPassword !== confirmPassword) {
      alert("Passwords do not match!");
      return;
    }
    try {
      const res = await fetch('http://localhost:8080/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: newUsername, password: newPassword })
      });

      if (res.ok) {
        alert("Success: Admin account created instantly!");
        newUsername = ''; newPassword = ''; confirmPassword = '';
        window.location.hash = ''; 
        currentView = 'login';
      } else {
        alert("Error: Username might already be taken.");
      }
    } catch (error) {
      console.error("Connection error:", error);
    }
  }

  async function handlePostAnnouncement() {
    if (!draftTitle || !draftMessage) {
      alert("Both a title and message are required.");
      return;
    }
    try {
      const res = await fetch('http://localhost:8080/api/announcements', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: draftTitle, message: draftMessage, author_id: username }) 
      });

      if (res.ok) {
        draftTitle = ''; draftMessage = '';
        await fetchAnnouncements(); 
      } else {
        alert("Server error: Could not save the announcement.");
      }
    } catch (error) {
      console.error("Posting error:", error);
    }
  }

  /** @param {string} id */
  async function deleteAnnouncement(id) {
    if (!confirm("Are you sure you want to delete this announcement forever?")) {
      return; 
    }

    try {
      const res = await fetch(`http://localhost:8080/api/announcements?id=${id}`, {
        method: 'DELETE'
      });

      if (res.ok) {
        await fetchAnnouncements();
      } else {
        alert("Failed to delete the post from the database.");
      }
    } catch (error) {
      console.error("Delete error:", error);
    }
  }

  function handleLogout() {
    currentView = 'public';
    username = ''; password = '';
  }
</script>

<main class="system-container">
  
  {#if currentView === 'public'}
    <div class="dashboard-view">
      <header class="public-header">
        <h1>Public Announcement Board</h1>
        <button class="action-btn outline" on:click={goToLogin}>Admin Login</button>
      </header>

      <section class="content">
        {#if announcements.length === 0}
          <div class="empty-state">
            <p>No announcements have been posted yet.</p>
          </div>
        {:else}
          <div class="announcement-list">
            {#each announcements as post (post.id)}
              <div class="card">
                <div class="card-header">
                  <h3>{post.title}</h3>
                  <span class="author-badge">Posted by {post.author_id}</span>
                </div>
                <p>{post.message}</p>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    </div>

  {:else if currentView === 'login'}
    <section class="login-gate">
      <div class="login-box">
        <h2>Admin Portal</h2>
        <div class="input-group">
          <label for="usernameInput">Username</label>
          <input id="usernameInput" type="text" bind:value={username} placeholder="Enter username" />
        </div>
        <div class="input-group">
          <label for="passwordInput">Password</label>
          <input id="passwordInput" type="password" bind:value={password} placeholder="Enter password" />
        </div>
        <div class="button-group">
          <button class="action-btn cancel" on:click={cancelLogin}>Back</button>
          <button class="action-btn primary" on:click={handleLogin}>Log In</button>
        </div>
      </div>
    </section>

  {:else if currentView === 'admin'}
    <div class="dashboard-view admin-mode">
      <header class="admin-header">
        <h1>Admin Control Center</h1>
        <button class="action-btn danger" on:click={handleLogout}>Logout</button>
      </header>

      <section class="content">
        <div class="admin-controls">
          <h2>Create New Announcement</h2>
          <div class="post-form">
             <input type="text" bind:value={draftTitle} placeholder="Announcement Title..." class="form-input" />
             <textarea bind:value={draftMessage} placeholder="Write your message here..." class="form-input textarea"></textarea>
             <button class="action-btn create-btn" on:click={handlePostAnnouncement}>Post to Board</button>
          </div>
        </div>
        
        <hr />
        
        <h2>Current Announcements</h2>
        {#if announcements.length === 0}
          <p class="empty-state">Database is empty.</p>
        {:else}
          <div class="announcement-list">
            {#each announcements as post (post.id)}
              <div class="card admin-card">
                <div class="card-header">
                  <div class="title-group">
                    <h3>{post.title}</h3>
                    <span class="author-badge">Posted by {post.author_id}</span>
                  </div>
                  <button class="delete-btn" on:click={() => deleteAnnouncement(post.id)}>Delete</button>
                </div>
                <p>{post.message}</p>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    </div>

  {:else if currentView === 'register'}
    <section class="login-gate">
      <div class="login-box" style="border-top: 4px solid #e67e22;">
        <h2>Create Admin</h2>
        <p style="text-align:center; font-size: 0.85rem; color: #6b7280; margin-top: -10px; margin-bottom: 20px;">Secret Registration Portal</p>
        
        <div class="input-group">
          <label for="newUsernameInput">New Username</label>
          <input id="newUsernameInput" type="text" bind:value={newUsername} placeholder="Choose username" />
        </div>
        <div class="input-group">
          <label for="newPasswordInput">New Password</label>
          <input id="newPasswordInput" type="password" bind:value={newPassword} placeholder="Create password" />
        </div>
        <div class="input-group">
          <label for="confirmPasswordInput">Confirm Password</label>
          <input id="confirmPasswordInput" type="password" bind:value={confirmPassword} placeholder="Confirm password" />
        </div>
        
        <div class="button-group">
          <button class="action-btn cancel" style="flex: 0.5" on:click={cancelLogin}>Back</button>
          <button class="action-btn primary" style="background-color: #e67e22; flex: 1" on:click={handleRegister}>
            Create Account
          </button>
        </div>
      </div>
    </section>
  {/if}

</main>

<style>
  /* ========================================= */
  /* PREMIUM SAAS UI POLISH                    */
  /* ========================================= */

  :global(body) { 
    margin: 0; 
    font-family: 'Inter', system-ui, -apple-system, sans-serif; 
    background-color: #f3f4f6; 
    color: #1f2937;
  }
  
  .system-container { 
    max-width: 900px; 
    margin: 0 auto; 
    padding: 0 1rem;
  }
  
  /* --- Premium Buttons --- */
  .action-btn { 
    padding: 0.6rem 1.2rem; 
    border-radius: 6px; 
    cursor: pointer; 
    font-weight: 600; 
    font-size: 0.9rem;
    border: none; 
    transition: all 0.2s ease-in-out;
  }
  .action-btn.primary { background-color: #4f46e5; color: white; box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.3); }
  .action-btn.primary:hover { background-color: #4338ca; transform: translateY(-1px); }
  
  .action-btn.outline { background-color: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.3); color: white; backdrop-filter: blur(4px); }
  .action-btn.outline:hover { background-color: rgba(255,255,255,0.2); }
  
  .action-btn.danger { background-color: #ef4444; color: white; }
  .action-btn.danger:hover { background-color: #dc2626; }
  
  .action-btn.cancel { background-color: #e5e7eb; color: #374151; }
  .action-btn.cancel:hover { background-color: #d1d5db; }
  
  .create-btn { background-color: #10b981; color: white; width: 100%; padding: 0.8rem; font-size: 1rem; margin-top: 0.5rem; box-shadow: 0 4px 6px -1px rgba(16, 185, 129, 0.3);}
  .create-btn:hover { background-color: #059669; transform: translateY(-1px); }
  
  .delete-btn { background-color: #fee2e2; color: #ef4444; border: 1px solid #fca5a5; padding: 0.3rem 0.8rem; border-radius: 6px; cursor: pointer; font-size: 0.8rem; font-weight: 600; transition: all 0.2s;}
  .delete-btn:hover { background-color: #ef4444; color: white; }

  /* --- Dashboards & Layout --- */
  .dashboard-view { 
    background: #f9fafb; 
    min-height: 100vh; 
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.05); 
  }
  
  header { 
    display: flex; 
    justify-content: space-between; 
    align-items: center; 
    padding: 1.5rem 2.5rem; 
    color: white; 
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }
  header h1 { margin: 0; font-size: 1.5rem; font-weight: 700; letter-spacing: -0.025em; white-space: nowrap; }
  
  .public-header { background: linear-gradient(135deg, #1e293b, #0f172a); }
  .admin-header { background: linear-gradient(135deg, #4f46e5, #312e81); } 
  
  .content { padding: 2.5rem; }
  
  /* --- The Login Gate --- */
  .login-gate { display: flex; justify-content: center; align-items: center; height: 100vh; background: #e5e7eb; }
  .login-box { 
    background: white; 
    padding: 2.5rem; 
    border-radius: 12px; 
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); 
    width: 340px; 
  }
  .login-box h2 { text-align: center; color: #111827; margin-top: 0; font-weight: 800; font-size: 1.5rem; margin-bottom: 1.5rem; }
  
  .input-group { margin-bottom: 1.2rem; }
  .input-group label { display: block; font-size: 0.85rem; font-weight: 600; margin-bottom: 0.4rem; color: #4b5563; }
  .input-group input { 
    width: 100%; 
    padding: 0.75rem; 
    border: 1px solid #d1d5db; 
    border-radius: 6px; 
    box-sizing: border-box; 
    transition: border-color 0.2s, box-shadow 0.2s;
  }
  .input-group input:focus { outline: none; border-color: #4f46e5; box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2); }
  .button-group { display: flex; justify-content: space-between; gap: 12px; margin-top: 2rem; }
  
  /* --- Posting Form --- */
  .post-form { 
    background: white; 
    border: 1px solid #e5e7eb; 
    padding: 1.5rem; 
    border-radius: 8px; 
    margin-bottom: 2.5rem; 
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  }
  .form-input { 
    width: 100%; 
    padding: 0.75rem; 
    margin-bottom: 1rem; 
    border: 1px solid #d1d5db; 
    border-radius: 6px; 
    box-sizing: border-box; 
    font-family: inherit; 
    font-size: 0.95rem;
    transition: border-color 0.2s, box-shadow 0.2s;
  }
  .form-input:focus { outline: none; border-color: #4f46e5; box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2); }
  .textarea { min-height: 120px; resize: vertical; line-height: 1.5;}
  
  /* --- Cards & Data --- */
  .empty-state { text-align: center; padding: 4rem 2rem; background: white; border-radius: 8px; border: 1px dashed #d1d5db; color: #6b7280; font-weight: 500;}
  
  .card { 
    border-left: 4px solid #3b82f6; 
    padding: 1.5rem; 
    background: white; 
    margin-bottom: 1.2rem; 
    border-radius: 0 8px 8px 0;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); 
    transition: transform 0.2s;
  }
  .card:hover { transform: translateY(-2px); box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); }
  .admin-card { border-left-color: #8b5cf6; }
  
  .card-header { display: flex; justify-content: space-between; align-items: flex-start; border-bottom: 1px solid #f3f4f6; padding-bottom: 0.75rem; margin-bottom: 0.75rem;}
  .title-group { display: flex; flex-direction: column; gap: 6px; }
  .title-group h3 { margin: 0; color: #1f2937; font-size: 1.2rem; font-weight: 700;}
  .author-badge { font-size: 0.75rem; background: #f3f4f6; padding: 0.25rem 0.6rem; border-radius: 9999px; color: #4b5563; font-weight: 600; align-self: flex-start;}
  
  .card p { margin: 0; color: #4b5563; line-height: 1.6; }
  hr { border: 0; border-top: 1px solid #e5e7eb; margin: 2.5rem 0; }
</style>