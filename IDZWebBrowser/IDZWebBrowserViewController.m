
#import "IDZWebBrowserViewController.h"
#import "_cgo_export.h"
#import "MBProgressHUD.h"
#import "notify.h"

@interface IDZWebBrowserViewController () <UIWebViewDelegate>
{
    BOOL   cancelled;
}
@property (weak, nonatomic) IBOutlet UIWebView *webView;
@property (weak, nonatomic) IBOutlet UIBarButtonItem *refresh;


- (void)loadRequestFromString:(NSString*)urlString;
- (void)updateButtons;
@end

@implementation IDZWebBrowserViewController

-(void)registerAppforDetectLockState {
    
    int notify_token;
    notify_register_dispatch("com.apple.springboard.lockstate",     &notify_token,dispatch_get_main_queue(), ^(int token) {
        uint64_t state = UINT64_MAX;
        notify_get_state(token, &state);
        if(state == 1) {
            NSLog(@"lock device");
            ShutDown();
        }
        
        NSLog(@"com.apple.springboard.lockstate = %llu", state);
        UILocalNotification *notification = [[UILocalNotification alloc]init];
        /*notification.repeatInterval = NSDayCalendarUnit;
        [notification setAlertBody:@"Hello world!! I come becoz you lock/unlock your device :)"];
        notification.alertAction = @"View";
        notification.alertAction = @"Yes";
        [notification setFireDate:[NSDate dateWithTimeIntervalSinceNow:1]];
        notification.soundName = UILocalNotificationDefaultSoundName;
        [notification setTimeZone:[NSTimeZone  defaultTimeZone]];*/
        
        [[UIApplication sharedApplication] presentLocalNotificationNow:notification];
        
    });
    
}


- (void)viewDidLoad
{
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(reloadWebView)
     
    name:UIApplicationWillEnterForegroundNotification object:nil];
    
    MBProgressHUD *hud = [MBProgressHUD showHUDAddedTo:self.view animated:YES];
    hud.labelText = @"Loading";
    
    [super viewDidLoad];
    [self registerAppforDetectLockState];
    
    dispatch_queue_t globalConcurrentQ = dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0);
    dispatch_async(globalConcurrentQ, ^{
        
        
        //while (!cancelled) {
            NSLog(@"--------Run Service------------");
            WebServer();
        //}
        //NSLog(@"---------Service Stopped---------");
        
    });

    [[UIApplication sharedApplication] setIdleTimerDisabled:YES];
    

    
    [self reloadWebView];

}

- (void)reloadWebView
{
    //MBProgressHUD *hud = [MBProgressHUD showHUDAddedTo:self.view animated:YES];
    //hud.labelText = @"Loading";

    //[UIApplication sharedApplication].networkActivityIndicatorVisible = YES;
    //[self updateButtons];

    NSAssert([self.webView isKindOfClass:[UIWebView class]], @"You webView outlet is not correctly connected.");
    self.webView.delegate = self;
    [self loadRequestFromString:@"http://127.0.0.1:8384"];
    
    //[UIApplication sharedApplication].networkActivityIndicatorVisible = NO;
    //[MBProgressHUD hideHUDForView:self.view animated:YES];
}


- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (void)loadRequestFromString:(NSString*)urlString
{
    NSURL *url = [NSURL URLWithString:urlString];
    NSURLRequest *urlRequest = [NSURLRequest requestWithURL:url];
    [self.webView loadRequest:urlRequest];
}

#pragma mark - Updating the UI
- (void)updateButtons
{
   
}

#pragma mark - UIWebViewDelegate
- (BOOL)webView:(UIWebView *)webView shouldStartLoadWithRequest:(NSURLRequest *)request navigationType:(UIWebViewNavigationType)navigationType
{
    return YES;
}


- (void)webViewDidStartLoad:(UIWebView *)webView
{
    [UIApplication sharedApplication].networkActivityIndicatorVisible = YES;
    [self updateButtons];
}
- (void)webViewDidFinishLoad:(UIWebView *)webView
{
    NSLog(@"%s", __PRETTY_FUNCTION__);
    [UIApplication sharedApplication].networkActivityIndicatorVisible = NO;
    [MBProgressHUD hideHUDForView:self.view animated:YES];
    [self updateButtons];
    
}
- (void)webView:(UIWebView *)webView didFailLoadWithError:(NSError *)error
{
    [UIApplication sharedApplication].networkActivityIndicatorVisible = NO;
    [self updateButtons];
    [self reloadWebView];
}



@end
